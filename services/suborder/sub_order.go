package services

import (
	"context"
	"order-service/clients"
	paymentClient "order-service/clients/payment"
	errorGeneral "order-service/constant/error"
	orderDTO "order-service/domain/dto/order"
	orderModel "order-service/domain/models/order"
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"

	"order-service/constant"
	errOrder "order-service/constant/error/order"
	orderHistoryDTO "order-service/domain/dto/orderhistory"
	orderPaymentDTO "order-service/domain/dto/orderpayment"
	subOrderDTO "order-service/domain/dto/suborder"
	subOrderModel "order-service/domain/models/suborder"
	"order-service/repositories"
	"order-service/utils/helper"
)

type ISubOrder struct {
	repository repositories.IRepositoryRegistry
	client     clients.IClientRegistry
}

type ISubOrderService interface {
	CreateOrder(context.Context, *subOrderDTO.SubOrderRequest) (*subOrderDTO.SubOrderResponse, error)
	Cancel(context.Context, string) error
	GetSubOrderList(context.Context, *subOrderDTO.SubOrderRequestParam) (*helper.PaginationResult, error)
	GetOrderDetail(context.Context, string) (*subOrderDTO.SubOrderResponse, error)
}

func NewSubOrderService(
	repository repositories.IRepositoryRegistry,
	client clients.IClientRegistry,
) ISubOrderService {
	return &ISubOrder{
		repository: repository,
		client:     client,
	}
}

func (o *ISubOrder) GetSubOrderList(
	ctx context.Context,
	request *subOrderDTO.SubOrderRequestParam,
) (*helper.PaginationResult, error) {
	var (
		orders []subOrderModel.SubOrder
		total  int64
	)
	orders, total, err := o.repository.GetSubOrderRepository().FindAllWithPagination(ctx, request)
	if err != nil {
		return nil, err
	}

	orderResponses := make([]subOrderDTO.SubOrderResponse, 0, len(orders))
	for _, order := range orders {
		orderResponses = append(orderResponses, subOrderDTO.SubOrderResponse{
			UUID:         order.UUID,
			SubOrderName: order.SubOrderName,
			CustomerID:   order.Order.CustomerID,
			PackageID:    order.Order.PackageID,
			Amount:       order.Amount,
			Status:       order.Status,
			IsPaid:       order.IsPaid,
			CompletedAt:  order.CompletedAt,
			CreatedAt:    order.CreatedAt,
			UpdatedAt:    order.UpdatedAt,
			Payment: &orderPaymentDTO.OrderPaymentResponse{
				PaymentID:   order.Payment.PaymentID,
				InvoiceID:   order.Payment.InvoiceID,
				PaymentLink: order.Payment.PaymentURL,
				Status:      order.Payment.Status,
			},
		})
	}

	pagination := helper.PaginationParam{
		Count: total,
		Page:  request.Page,
		Limit: request.Limit,
		Data:  orderResponses,
	}
	response := helper.GeneratePagination(pagination)
	return &response, nil
}

func (o *ISubOrder) GetOrderDetail(ctx context.Context, subOrderUUID string) (*subOrderDTO.SubOrderResponse, error) {
	var (
		subOrder *subOrderModel.SubOrder
	)

	subOrder, err := o.repository.GetSubOrderRepository().FindOneByUUID(ctx, subOrderUUID)
	if err != nil {
		return nil, err
	}

	response := &subOrderDTO.SubOrderResponse{
		UUID:         subOrder.UUID,
		SubOrderName: subOrder.SubOrderName,
		CustomerID:   subOrder.Order.CustomerID,
		PackageID:    subOrder.Order.PackageID,
		Amount:       subOrder.Amount,
		Status:       subOrder.Status,
		IsPaid:       subOrder.IsPaid,
		CompletedAt:  subOrder.CompletedAt,
		Payment: &orderPaymentDTO.OrderPaymentResponse{
			PaymentID:   subOrder.Payment.PaymentID,
			InvoiceID:   subOrder.Payment.InvoiceID,
			PaymentLink: subOrder.Payment.PaymentURL,
			Status:      subOrder.Payment.Status,
		},
	}
	return response, nil
}

//nolint:cyclop
func (o *ISubOrder) CreateOrder(
	ctx context.Context,
	request *subOrderDTO.SubOrderRequest,
) (*subOrderDTO.SubOrderResponse, error) {
	var (
		subOrder        *subOrderModel.SubOrder
		order           *orderModel.Order
		txErr           error
		completedAt     *time.Time
		paymentResponse *paymentClient.PaymentData
	)

	today := time.Now()
	if today.After(request.OrderDate) {
		return nil, errorGeneral.ErrOrderDate
	}

	tx := o.repository.GetTx()
	err := tx.Transaction(func(tx *gorm.DB) error {
		customerID, _ := uuid.Parse(request.CustomerID) //nolint:errcheck
		order, txErr = o.repository.GetOrder().FindOneOrderByCustomerIDWithLocking(ctx, customerID)
		if txErr != nil {
			return txErr
		}

		if order != nil {
			return errOrder.ErrOrderNotEmpty
		}

		if request.PaymentType == constant.PTFullPayment {
			utcTimeAddSevenHour := time.Now().UTC().Add(7 * time.Hour)
			completedAt = &utcTimeAddSevenHour
		}

		order, txErr = o.repository.GetOrder().Create(ctx, tx, &orderDTO.OrderRequest{
			CustomerID:  request.CustomerID,
			PackageID:   request.PackageID,
			OrderDate:   request.OrderDate,
			CompletedAt: completedAt,
		})
		if txErr != nil {
			return txErr
		}

		subOrder, txErr = o.repository.GetSubOrderRepository().Create(ctx, tx, &subOrderDTO.SubOrderRequest{
			OrderID:     order.ID,
			Status:      constant.Pending,
			Amount:      request.Amount,
			PaymentType: request.PaymentType,
		})
		if txErr != nil {
			return txErr
		}

		txErr = o.repository.GetOrderHistoryRepository().Create(ctx, tx, &orderHistoryDTO.OrderHistoryRequest{
			SubOrderID: subOrder.ID,
			Status:     constant.PendingString,
		})
		if txErr != nil {
			return txErr
		}

		paymentResponse, txErr = o.client.GetPayment().CreatePaymentLink(&paymentClient.PaymentRequest{
			OrderID:     subOrder.UUID,
			ExpiredAt:   time.Now().Add(24 * time.Hour),
			Amount:      request.Amount,
			Description: request.PaymentType.Title(),
			CustomerDetail: paymentClient.CustomerDetail{
				Name:  "Customer",
				Email: "customer@test.com",
				Phone: "081234567890",
			},
			ItemDetail: []paymentClient.ItemDetail{
				{
					ID:       uuid.New(),
					Name:     request.PaymentType.Title(),
					Amount:   request.Amount,
					Quantity: 1,
				},
			},
		})
		if txErr != nil {
			return txErr
		}

		txErr = o.repository.GetOrderPaymentRepository().
			Create(ctx, tx, &orderPaymentDTO.OrderPaymentRequest{
				SubOrderID:  subOrder.ID,
				PaymentID:   paymentResponse.UUID,
				PaymentLink: paymentResponse.PaymentLink,
				InvoiceID:   uuid.New(),
				Status:      paymentResponse.Status,
			})
		if txErr != nil {
			return txErr
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	response := subOrderDTO.SubOrderResponse{
		UUID:         subOrder.UUID,
		SubOrderName: subOrder.SubOrderName,
		CustomerID:   order.CustomerID,
		PackageID:    order.PackageID,
		Amount:       subOrder.Amount,
		Status:       subOrder.Status,
		OrderDate:    order.OrderDate,
		IsPaid:       subOrder.IsPaid,
		Payment: &orderPaymentDTO.OrderPaymentResponse{
			PaymentID:   paymentResponse.UUID,
			InvoiceID:   uuid.New(),
			PaymentLink: paymentResponse.PaymentLink,
			Status:      paymentResponse.Status,
		},
	}
	return &response, nil
}

func (o *ISubOrder) Cancel(ctx context.Context, subOrderUUID string) error {
	var (
		order *subOrderModel.SubOrder
		txErr error
	)

	tx := o.repository.GetTx()
	err := tx.Transaction(func(tx *gorm.DB) error {
		order, txErr = o.repository.GetSubOrderRepository().FindOneByUUID(ctx, subOrderUUID)
		if txErr != nil {
			return txErr
		}

		if order != nil && order.Status == constant.Cancelled {
			return errOrder.ErrCancelOrder
		}

		uuid, _ := uuid.Parse(subOrderUUID) //nolint:errcheck
		txErr = o.repository.GetSubOrderRepository().Cancel(ctx, tx, &subOrderDTO.CancelRequest{
			UUID:   uuid,
			Status: constant.Cancelled,
		}, &subOrderModel.SubOrder{
			Status: order.Status,
		})
		if txErr != nil {
			return txErr
		}

		txErr = o.repository.GetOrderHistoryRepository().Create(ctx, tx, &orderHistoryDTO.OrderHistoryRequest{
			SubOrderID: order.ID,
			Status:     constant.CancelledString,
		})
		if txErr != nil {
			return txErr
		}

		txErr = o.repository.GetOrder().DeleteByOrderID(ctx, order.OrderID)
		if txErr != nil {
			return txErr
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}