package services

import (
	"context"
	"time"

	"gorm.io/gorm"

	"order-service/clients"
	paymentClient "order-service/clients/payment"
	rbacClient "order-service/clients/rbac"
	errorGeneral "order-service/constant/error"
	orderDTO "order-service/domain/dto/order"
	"order-service/utils/sentry"

	"github.com/google/uuid"

	"order-service/constant"
	errOrder "order-service/constant/error/order"
	orderHistoryDTO "order-service/domain/dto/orderhistory"
	orderPaymentDTO "order-service/domain/dto/orderpayment"
	subOrderDTO "order-service/domain/dto/suborder"
	"order-service/domain/models"
	"order-service/repositories"
	"order-service/utils/helper"
)

type ISubOrder struct {
	repository repositories.IRepositoryRegistry
	client     clients.IClientRegistry
	sentry     sentry.ISentry
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
	sentry sentry.ISentry,
) ISubOrderService {
	return &ISubOrder{
		repository: repository,
		client:     client,
		sentry:     sentry,
	}
}

func (o *ISubOrder) GetSubOrderList(
	ctx context.Context,
	request *subOrderDTO.SubOrderRequestParam,
) (*helper.PaginationResult, error) {
	const logCtx = "services.suborder.sub_order.GetSubOrderList"
	var (
		orders []models.SubOrder
		total  int64
		span   = o.sentry.StartSpan(ctx, logCtx)
	)
	ctx = o.sentry.SpanContext(span)
	defer o.sentry.Finish(span)

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
	const logCtx = "services.suborder.sub_order.GetOrderDetail"
	var (
		subOrder *models.SubOrder
		span     = o.sentry.StartSpan(ctx, logCtx)
	)
	ctx = o.sentry.SpanContext(span)
	defer o.sentry.Finish(span)

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
	const logCtx = "services.suborder.sub_order.CreateOrder"
	var (
		response *subOrderDTO.SubOrderResponse
		err      error
		span     = o.sentry.StartSpan(ctx, logCtx)
	)
	ctx = o.sentry.SpanContext(span)
	defer o.sentry.Finish(span)

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
	const logCtx = "services.suborder.sub_order.createDownPaymentOrder"
	var (
		subOrder         *models.SubOrder
		order            *models.Order
		txErr            error
		paymentResponse  *paymentClient.PaymentData
		customerResponse *rbacClient.RBACData
		err              error
		orderHistories   []orderHistoryDTO.OrderHistoryRequest
		amount           float64 = 100000000 // will be remove if package service is ready
		span                     = o.sentry.StartSpan(ctx, logCtx)
	)
	ctx = o.sentry.SpanContext(span)
	defer o.sentry.Finish(span)

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
		order, txErr = o.repository.GetOrder().FindOneOrderByCustomerIDWithLocking(ctx, tx, customerID)
		if txErr != nil {
			return txErr
		}

		if order != nil {
			if len(order.SubOrder) != 3 || order.CompletedAt == nil {
				return errOrder.ErrPreviousOrderNotEmpty
			}
		}

		order, txErr = o.repository.GetOrder().Create(ctx, tx, &orderDTO.OrderRequest{
			CustomerID:                 request.CustomerID.String(),
			PackageID:                  request.PackageID.String(),
			RemainingOutstandingAmount: amount,
		})
		if txErr != nil {
			return txErr
		}

		subOrder, txErr = o.repository.GetSubOrder().Create(ctx, tx, &models.SubOrder{
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
	const logCtx = "services.suborder.sub_order.createHalfPaymentOrder"
	var (
		subOrder         *models.SubOrder
		order            *models.Order
		txErr            error
		paymentResponse  *paymentClient.PaymentData
		customerResponse *rbacClient.RBACData
		err              error
		orderHistories   []orderHistoryDTO.OrderHistoryRequest
		span             = o.sentry.StartSpan(ctx, logCtx)
	)
	ctx = o.sentry.SpanContext(span)
	defer o.sentry.Finish(span)

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

	if subOrder != nil {
		if *subOrder.IsPaid {
			return nil, errOrder.ErrHalfPaymentNotEmpty
		} else { //nolint:revive
			return nil, errOrder.ErrPreviousOrderNotEmpty
		}
	}

	total := order.RemainingOutstandingAmount * 50 / 100
	if total != request.Amount {
		return nil, errOrder.ErrInvalidHalfAmount
	}
	tx := o.repository.GetTx()
	err = tx.Transaction(func(tx *gorm.DB) error { //nolint:dupl
		subOrder, txErr = o.repository.GetSubOrder().Create(ctx, tx, &models.SubOrder{
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
	const logCtx = "services.suborder.sub_order.createFullPaymentOrder"
	var (
		subOrder         *models.SubOrder
		order            *models.Order
		txErr            error
		paymentResponse  *paymentClient.PaymentData
		customerResponse *rbacClient.RBACData
		err              error
		orderHistories   []orderHistoryDTO.OrderHistoryRequest
		span             = o.sentry.StartSpan(ctx, logCtx)
	)
	ctx = o.sentry.SpanContext(span)
	defer o.sentry.Finish(span)

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
			constant.PTFullPayment.String())
	if err != nil {
		return nil, err
	}

	if subOrder != nil {
		if *subOrder.IsPaid {
			return nil, errOrder.ErrFullPaymentNotEmpty
		} else { //nolint:revive
			return nil, errOrder.ErrPreviousOrderNotEmpty
		}
	}

	total := order.RemainingOutstandingAmount - request.Amount
	if total != 0 {
		return nil, errOrder.ErrInvalidFullAmount
	}
	tx := o.repository.GetTx()
	err = tx.Transaction(func(tx *gorm.DB) error { //nolint:dupl
		subOrder, txErr = o.repository.GetSubOrder().Create(ctx, tx, &models.SubOrder{
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
	const logCtx = "services.suborder.sub_order.Cancel"
	var (
		order *models.SubOrder
		txErr error
		span  = o.sentry.StartSpan(ctx, logCtx)
	)
	ctx = o.sentry.SpanContext(span)
	defer o.sentry.Finish(span)

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
		}, &models.SubOrder{
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
	const logCtx = "services.suborder.sub_order.processPayment"
	var (
		updateRequest       subOrderDTO.UpdateSubOrderRequest
		paidAt, completedAt *time.Time
		isPaid              = false
		order               *models.Order
		total               float64
		span                = o.sentry.StartSpan(ctx, logCtx)
	)
	ctx = o.sentry.SpanContext(span)
	defer o.sentry.Finish(span)

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
			completedAtTime := time.Now()
			completedAt = &completedAtTime
			total = order.RemainingOutstandingAmount - request.Amount
			updateRequest = subOrderDTO.UpdateSubOrderRequest{
				Status: constant.PaymentSuccess,
				IsPaid: &isPaid,
			}
		case constant.Cancelled:
			canceledAt := time.Now()
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

		txErr := o.repository.GetSubOrder().Update(ctx, tx, &updateRequest, &models.SubOrder{
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
			updateOrder := &orderDTO.OrderRequest{
				OrderID:                    order.UUID.String(),
				RemainingOutstandingAmount: total,
			}

			if subOrder.PaymentType == constant.PTFullPayment {
				updateOrder.CompletedAt = completedAt
			}

			txErr = o.repository.GetOrder().Update(ctx, tx, updateOrder)
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
