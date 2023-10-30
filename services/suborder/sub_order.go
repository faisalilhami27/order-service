package services

import (
	"context"
	"time"

	"order-service/clients"
	paymentClient "order-service/clients/payment"
	rbacClient "order-service/clients/rbac"
	errorGeneral "order-service/constant/error"
	orderDTO "order-service/domain/dto/order"
	orderModel "order-service/domain/models/order"

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
	ReceivePendingPayment(context.Context, *subOrderDTO.PaymentRequest) error
	ReceivePaymentSettlement(context.Context, *subOrderDTO.PaymentRequest) error
	ReceivePaymentExpire(context.Context, *subOrderDTO.PaymentRequest) error
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
	orders, total, err := o.repository.GetSubOrder().FindAllWithPagination(ctx, request)
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
			CreatedAt:    order.CreatedAt,
			UpdatedAt:    order.UpdatedAt,
			Payment: &orderPaymentDTO.OrderPaymentResponse{
				PaymentID:   order.Payment.PaymentID,
				InvoiceID:   order.Payment.InvoiceID,
				PaymentLink: *order.Payment.PaymentURL,
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

	subOrder, err := o.repository.GetSubOrder().FindOneByUUID(ctx, subOrderUUID)
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
		Payment: &orderPaymentDTO.OrderPaymentResponse{
			PaymentID:   subOrder.Payment.PaymentID,
			InvoiceID:   subOrder.Payment.InvoiceID,
			PaymentLink: *subOrder.Payment.PaymentURL,
			Status:      subOrder.Payment.Status,
		},
	}
	return response, nil
}

func (o *ISubOrder) CreateOrder(
	ctx context.Context,
	request *subOrderDTO.SubOrderRequest,
) (*subOrderDTO.SubOrderResponse, error) {
	var (
		response *subOrderDTO.SubOrderResponse
		err      error
	)
	switch request.PaymentType {
	case constant.PTDownPayment:
		response, err = o.createDownPaymentOrder(ctx, request)
	case constant.PTHalfPayment:
		response, err = o.createHalfPaymentOrder(ctx, request)
	case constant.PTFullPayment:
		response, err = o.createFullPaymentOrder(ctx, request)
	}
	if err != nil {
		return nil, err
	}

	return response, err
}

//nolint:cyclop
func (o *ISubOrder) createDownPaymentOrder(
	ctx context.Context,
	request *subOrderDTO.SubOrderRequest,
) (*subOrderDTO.SubOrderResponse, error) {
	var (
		subOrder         *subOrderModel.SubOrder
		order            *orderModel.Order
		txErr            error
		paymentResponse  *paymentClient.PaymentData
		customerResponse *rbacClient.RBACData
		err              error
		orderHistories   []orderHistoryDTO.OrderHistoryRequest
		amount           float64 = 100000000 // will be remove if package service is ready
	)

	today := time.Now()
	if today.After(request.OrderDate) {
		return nil, errorGeneral.ErrOrderDate
	}

	total := amount * 10 / 100
	if total != request.Amount {
		return nil, errOrder.ErrInvalidDownAmount
	}
	tx := o.repository.GetTx()
	err = tx.Transaction(func(tx *gorm.DB) error {
		customerID, _ := uuid.Parse(request.CustomerID.String()) //nolint:errcheck
		order, txErr = o.repository.GetOrder().FindOneOrderByCustomerIDWithLocking(ctx, customerID)
		if txErr != nil {
			return txErr
		}

		if order != nil {
			return errOrder.ErrPreviousOrderNotEmpty
		}

		order, txErr = o.repository.GetOrder().Create(ctx, tx, &orderDTO.OrderRequest{
			CustomerID:                 request.CustomerID.String(),
			PackageID:                  request.PackageID.String(),
			RemainingOutstandingAmount: amount,
		})
		if txErr != nil {
			return txErr
		}

		subOrder, txErr = o.repository.GetSubOrder().Create(ctx, tx, &subOrderModel.SubOrder{
			OrderID:     order.ID,
			Status:      constant.Pending,
			Amount:      request.Amount,
			PaymentType: request.PaymentType,
			OrderDate:   request.OrderDate,
		})
		if txErr != nil {
			return txErr
		}

		orderHistories = []orderHistoryDTO.OrderHistoryRequest{
			{
				SubOrderID: subOrder.ID,
				Status:     constant.PendingString,
			},
		}
		txErr = o.repository.GetOrderHistory().BulkCreate(ctx, tx, orderHistories)
		if txErr != nil {
			return txErr
		}

		customerResponse, txErr = o.client.GetRBAC().GetUserRBAC(request.CustomerID.String())
		if txErr != nil {
			return txErr
		}

		expiredAt := time.Now().Add(24 * time.Hour)
		paymentResponse, txErr = o.client.GetPayment().CreatePaymentLink(&paymentClient.PaymentRequest{
			OrderID:     subOrder.UUID,
			ExpiredAt:   expiredAt,
			Amount:      request.Amount,
			Description: request.PaymentType.Title(),
			CustomerDetail: paymentClient.CustomerDetail{
				Name:  customerResponse.Name,
				Email: customerResponse.Email,
				Phone: customerResponse.PhoneNumber,
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

		txErr = o.repository.GetOrderPayment().
			Create(ctx, tx, &orderPaymentDTO.OrderPaymentRequest{
				Amount:      request.Amount,
				SubOrderID:  subOrder.ID,
				PaymentID:   paymentResponse.UUID,
				PaymentLink: paymentResponse.PaymentLink,
				InvoiceID:   uuid.New(),
				Status:      paymentResponse.Status,
				ExpiredAt:   &expiredAt,
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
		OrderDate:    subOrder.OrderDate,
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

//nolint:cyclop
func (o *ISubOrder) createHalfPaymentOrder(
	ctx context.Context,
	request *subOrderDTO.SubOrderRequest,
) (*subOrderDTO.SubOrderResponse, error) {
	var (
		subOrder         *subOrderModel.SubOrder
		order            *orderModel.Order
		txErr            error
		paymentResponse  *paymentClient.PaymentData
		customerResponse *rbacClient.RBACData
		err              error
		orderHistories   []orderHistoryDTO.OrderHistoryRequest
	)

	today := time.Now()
	if today.After(request.OrderDate) {
		return nil, errorGeneral.ErrOrderDate
	}

	order, err = o.repository.GetOrder().FindOneOrderByUUID(ctx, request.OrderID)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, errOrder.ErrOrderIsEmpty
	}

	subOrder, err = o.repository.GetSubOrder().
		FindOneByOrderIDAndPaymentType(
			ctx,
			order.ID,
			constant.PTDownPayment.String())
	if err != nil {
		return nil, err
	}

	if subOrder == nil {
		return nil, errOrder.ErrHalfPaymentIsEmpty
	}

	total := order.RemainingOutstandingAmount * 50 / 100
	if total != request.Amount {
		return nil, errOrder.ErrInvalidHalfAmount
	}
	tx := o.repository.GetTx()
	err = tx.Transaction(func(tx *gorm.DB) error { //nolint:dupl
		subOrder, txErr = o.repository.GetSubOrder().Create(ctx, tx, &subOrderModel.SubOrder{
			OrderID:     order.ID,
			Status:      constant.Pending,
			Amount:      request.Amount,
			PaymentType: request.PaymentType,
			OrderDate:   request.OrderDate,
		})
		if txErr != nil {
			return txErr
		}

		orderHistories = []orderHistoryDTO.OrderHistoryRequest{
			{
				SubOrderID: subOrder.ID,
				Status:     constant.PendingString,
			},
		}
		txErr = o.repository.GetOrderHistory().BulkCreate(ctx, tx, orderHistories)
		if txErr != nil {
			return txErr
		}

		customerResponse, txErr = o.client.GetRBAC().GetUserRBAC(request.CustomerID.String())
		if txErr != nil {
			return txErr
		}

		expiredAt := time.Now().Add(24 * time.Hour)
		paymentResponse, txErr = o.client.GetPayment().CreatePaymentLink(&paymentClient.PaymentRequest{
			OrderID:     subOrder.UUID,
			ExpiredAt:   expiredAt,
			Amount:      request.Amount,
			Description: request.PaymentType.Title(),
			CustomerDetail: paymentClient.CustomerDetail{
				Name:  customerResponse.Name,
				Email: customerResponse.Email,
				Phone: customerResponse.PhoneNumber,
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

		txErr = o.repository.GetOrderPayment().
			Create(ctx, tx, &orderPaymentDTO.OrderPaymentRequest{
				Amount:      request.Amount,
				SubOrderID:  subOrder.ID,
				PaymentID:   paymentResponse.UUID,
				PaymentLink: paymentResponse.PaymentLink,
				InvoiceID:   uuid.New(),
				Status:      paymentResponse.Status,
				ExpiredAt:   &expiredAt,
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
		OrderDate:    subOrder.OrderDate,
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

//nolint:cyclop
func (o *ISubOrder) createFullPaymentOrder(
	ctx context.Context,
	request *subOrderDTO.SubOrderRequest,
) (*subOrderDTO.SubOrderResponse, error) {
	var (
		subOrder         *subOrderModel.SubOrder
		order            *orderModel.Order
		txErr            error
		paymentResponse  *paymentClient.PaymentData
		customerResponse *rbacClient.RBACData
		err              error
		orderHistories   []orderHistoryDTO.OrderHistoryRequest
	)

	today := time.Now()
	if today.After(request.OrderDate) {
		return nil, errorGeneral.ErrOrderDate
	}

	order, err = o.repository.GetOrder().FindOneOrderByUUID(ctx, request.OrderID)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, errOrder.ErrOrderIsEmpty
	}

	subOrder, err = o.repository.GetSubOrder().
		FindOneByOrderIDAndPaymentType(
			ctx,
			order.ID,
			constant.PTHalfPayment.String())
	if err != nil {
		return nil, err
	}

	if subOrder == nil {
		return nil, errOrder.ErrInvalidFullAmount
	}

	total := order.RemainingOutstandingAmount - request.Amount
	if total != 0 {
		return nil, errOrder.ErrInvalidFullAmount
	}
	tx := o.repository.GetTx()
	err = tx.Transaction(func(tx *gorm.DB) error { //nolint:dupl
		subOrder, txErr = o.repository.GetSubOrder().Create(ctx, tx, &subOrderModel.SubOrder{
			OrderID:     order.ID,
			Status:      constant.Pending,
			Amount:      request.Amount,
			PaymentType: request.PaymentType,
			OrderDate:   request.OrderDate,
		})
		if txErr != nil {
			return txErr
		}

		orderHistories = []orderHistoryDTO.OrderHistoryRequest{
			{
				SubOrderID: subOrder.ID,
				Status:     constant.PendingString,
			},
		}
		txErr = o.repository.GetOrderHistory().BulkCreate(ctx, tx, orderHistories)
		if txErr != nil {
			return txErr
		}

		customerResponse, txErr = o.client.GetRBAC().GetUserRBAC(request.CustomerID.String())
		if txErr != nil {
			return txErr
		}

		expiredAt := time.Now().Add(24 * time.Hour)
		paymentResponse, txErr = o.client.GetPayment().CreatePaymentLink(&paymentClient.PaymentRequest{
			OrderID:     subOrder.UUID,
			ExpiredAt:   expiredAt,
			Amount:      request.Amount,
			Description: request.PaymentType.Title(),
			CustomerDetail: paymentClient.CustomerDetail{
				Name:  customerResponse.Name,
				Email: customerResponse.Email,
				Phone: customerResponse.PhoneNumber,
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

		txErr = o.repository.GetOrderPayment().
			Create(ctx, tx, &orderPaymentDTO.OrderPaymentRequest{
				Amount:      request.Amount,
				SubOrderID:  subOrder.ID,
				PaymentID:   paymentResponse.UUID,
				PaymentLink: paymentResponse.PaymentLink,
				InvoiceID:   uuid.New(),
				Status:      paymentResponse.Status,
				ExpiredAt:   &expiredAt,
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
		OrderDate:    subOrder.OrderDate,
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
		order, txErr = o.repository.GetSubOrder().FindOneByUUID(ctx, subOrderUUID)
		if txErr != nil {
			return txErr
		}

		if order != nil && order.Status == constant.Cancelled {
			return errOrder.ErrCancelOrder
		}

		uuidParse, _ := uuid.Parse(subOrderUUID) //nolint:errcheck
		txErr = o.repository.GetSubOrder().Cancel(ctx, tx, &subOrderDTO.CancelRequest{
			UUID:   uuidParse,
			Status: constant.Cancelled,
		}, &subOrderModel.SubOrder{
			Status: order.Status,
		})
		if txErr != nil {
			return txErr
		}

		txErr = o.repository.GetOrderHistory().Create(ctx, tx, &orderHistoryDTO.OrderHistoryRequest{
			SubOrderID: order.ID,
			Status:     constant.CancelledString,
		})
		if txErr != nil {
			return txErr
		}

		txErr = o.repository.GetOrder().DeleteByOrderID(ctx, tx, order.OrderID)
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

//nolint:cyclop
func (o *ISubOrder) processPayment(
	ctx context.Context,
	request *subOrderDTO.PaymentRequest,
	status constant.OrderStatus,
) error {
	var (
		updateRequest       subOrderDTO.UpdateSubOrderRequest
		paidAt, completedAt *time.Time
		isPaid              = false
		order               *orderModel.Order
		total               float64
	)

	subOrder, err := o.repository.GetSubOrder().FindOneByUUID(ctx, request.OrderID.String())
	if err != nil {
		return err
	}

	order, err = o.repository.GetOrder().FindOneOrderByID(ctx, subOrder.OrderID)
	if err != nil {
		return err
	}

	tx := o.repository.GetTx()
	err = tx.Transaction(func(tx *gorm.DB) error {
		switch status {
		case constant.PaymentSuccess:
			isPaid = true
			paidAt = request.PaidAt
			total = order.RemainingOutstandingAmount - request.Amount
			updateRequest = subOrderDTO.UpdateSubOrderRequest{
				Status: constant.PaymentSuccess,
				IsPaid: &isPaid,
			}
		case constant.Cancelled:
			canceledAt := time.Now().UTC().Add(7 * time.Hour)
			updateRequest = subOrderDTO.UpdateSubOrderRequest{
				Status:     constant.Cancelled,
				CanceledAt: &canceledAt,
			}
		case constant.PendingPayment:
			updateRequest = subOrderDTO.UpdateSubOrderRequest{
				Status: constant.PendingPayment,
			}
		default:
			return errorGeneral.ErrStatus
		}

		txErr := o.repository.GetSubOrder().Update(ctx, tx, &updateRequest, &subOrderModel.SubOrder{
			UUID:   subOrder.UUID,
			Status: subOrder.Status,
		})
		if txErr != nil {
			return txErr
		}

		txErr = o.repository.GetOrderHistory().Create(ctx, tx, &orderHistoryDTO.OrderHistoryRequest{
			SubOrderID: subOrder.ID,
			Status:     constant.OrderStatusString(status.String()),
		})
		if txErr != nil {
			return txErr
		}

		txErr = o.repository.GetOrderPayment().
			Update(ctx, tx, &orderPaymentDTO.OrderPaymentRequest{
				Amount:      request.Amount,
				PaymentID:   request.PaymentID,
				PaymentLink: request.PaymentLink,
				PaymentType: &request.PaymentType,
				VANumber:    request.VaNumber,
				Bank:        request.Bank,
				Acquirer:    request.Acquirer,
				Status:      &request.Status,
				PaidAt:      paidAt,
			})
		if txErr != nil {
			return txErr
		}

		if request.Status == constant.PaymentStatusSettlement.String() {
			txErr = o.repository.GetOrder().Update(ctx, tx, &orderDTO.OrderRequest{
				OrderID:                    order.UUID.String(),
				RemainingOutstandingAmount: total,
				CompletedAt:                completedAt,
			})
			if txErr != nil {
				return txErr
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (o *ISubOrder) ReceivePendingPayment(ctx context.Context, request *subOrderDTO.PaymentRequest) error {
	return o.processPayment(ctx, request, constant.PendingPayment)
}

func (o *ISubOrder) ReceivePaymentSettlement(ctx context.Context, request *subOrderDTO.PaymentRequest) error {
	return o.processPayment(ctx, request, constant.PaymentSuccess)
}

func (o *ISubOrder) ReceivePaymentExpire(ctx context.Context, request *subOrderDTO.PaymentRequest) error {
	return o.processPayment(ctx, request, constant.Cancelled)
}
