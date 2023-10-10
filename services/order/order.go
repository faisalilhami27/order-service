package order

import (
	"context"
	orderPaymentDTO "order-service/domain/dto/orderpayment"
	orderPaymentModel "order-service/domain/models/orderpayment"
	"order-service/utils"

	"gorm.io/gorm"

	"github.com/google/uuid"

	"order-service/constant"
	errOrder "order-service/constant/error/order"
	orderDTO "order-service/domain/dto/order"
	orderHistoryDTO "order-service/domain/dto/orderhistory"
	orderModel "order-service/domain/models/order"
	"order-service/repositories"
)

type IOrder struct {
	repository repositories.IRepositoryRegistry
}

type IOrderService interface {
	CreateOrder(context.Context, *orderDTO.OrderRequest) (*orderDTO.OrderResponse, error)
	GetOrderList(context.Context, *orderDTO.OrderRequestParam) (*utils.PaginationResult, error)
}

func NewOrderService(repository repositories.IRepositoryRegistry) IOrderService {
	return &IOrder{
		repository: repository,
	}
}

func (o *IOrder) GetOrderList(
	ctx context.Context,
	request *orderDTO.OrderRequestParam,
) (*utils.PaginationResult, error) {
	var (
		orders []orderModel.Order
		total  int64
	)
	orders, total, err := o.repository.GetOrderRepository().FindAllWithPagination(ctx, request)
	if err != nil {
		return nil, err
	}

	orderResponses := make([]orderDTO.OrderResponse, 0, len(orders))
	for _, order := range orders {
		orderResponses = append(orderResponses, orderDTO.OrderResponse{
			UUID:        order.UUID,
			OrderName:   order.OrderName,
			CustomerID:  order.CustomerID,
			PackageID:   order.PackageID,
			Amount:      order.Amount,
			Status:      order.Status,
			OrderDate:   order.OrderDate,
			IsPaid:      order.IsPaid,
			CompletedAt: order.CompletedAt,
			CreatedAt:   order.CreatedAt,
			UpdatedAt:   order.UpdatedAt,
			Payment: &orderPaymentDTO.OrderPaymentResponse{
				PaymentID:  order.Payment.PaymentID,
				InvoiceID:  order.Payment.InvoiceID,
				PaymentURL: order.Payment.PaymentURL,
				Status:     order.Payment.Status,
			},
		})
	}

	pagination := utils.PaginationParam{
		Count: total,
		Page:  request.Page,
		Limit: request.Limit,
		Data:  orderResponses,
	}
	response := utils.GeneratePagination(pagination)
	return &response, nil
}

func (o *IOrder) CreateOrder(ctx context.Context, request *orderDTO.OrderRequest) (*orderDTO.OrderResponse, error) {
	var (
		orderResult, order *orderModel.Order
		orderPayment       *orderPaymentModel.OrderPayment
		txErr              error
	)

	tx := o.repository.GetTx()
	err := tx.Transaction(func(tx *gorm.DB) error {
		customerID, _ := uuid.Parse(request.CustomerID) //nolint:errcheck
		order, txErr = o.repository.GetOrderRepository().FindOneOrderByCustomerIDWithLocking(ctx, customerID)
		if txErr != nil {
			return txErr
		}

		if order != nil {
			return errOrder.ErrOrderNotEmpty
		}

		request.Status = constant.Pending
		orderResult, txErr = o.repository.GetOrderRepository().Create(ctx, tx, request)
		if txErr != nil {
			return txErr
		}

		txErr = o.repository.GetOrderHistoryRepository().Create(ctx, tx, &orderHistoryDTO.OrderHistoryRequest{
			OrderID: orderResult.ID,
			Status:  constant.PendingString,
		})
		if txErr != nil {
			return txErr
		}

		// Should! will continue if invoice and payment service is ready
		status := "pending"
		orderPayment, txErr = o.repository.GetOrderPaymentRepository().
			Create(ctx, tx, &orderPaymentDTO.OrderPaymentRequest{
				OrderID:    orderResult.ID,
				PaymentID:  uuid.New(),
				InvoiceID:  uuid.New(),
				PaymentURL: "https://payment.com",
				Status:     &status,
			})
		if txErr != nil {
			return txErr
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	response := orderDTO.ResponseFormatter(orderResult, &orderPaymentDTO.OrderPaymentResponse{
		PaymentID:  orderPayment.PaymentID,
		InvoiceID:  orderPayment.InvoiceID,
		PaymentURL: orderPayment.PaymentURL,
		Status:     orderPayment.Status,
	})
	return response, nil
}
