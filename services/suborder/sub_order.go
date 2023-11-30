package services

import (
	"context"
	"fmt"
	"math/rand"
	"order-service/config"

	invoiceModel "order-service/clients/invoice"
	packageClient "order-service/clients/package"

	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"

	"order-service/clients"
	paymentClient "order-service/clients/payment"
	rbacClient "order-service/clients/rbac"
	"order-service/common/circuitbreaker"
	"order-service/common/sentry"
	errorGeneral "order-service/constant/error"
	orderDTO "order-service/domain/dto/order"

	"order-service/constant"
	errOrder "order-service/constant/error/order"
	orderHistoryDTO "order-service/domain/dto/orderhistory"
	orderPaymentDTO "order-service/domain/dto/orderpayment"
	subOrderDTO "order-service/domain/dto/suborder"
	"order-service/domain/models"
	"order-service/repositories"
	"order-service/utils/helper"
)

type SubOrder struct {
	repository repositories.IRepositoryRegistry
	client     clients.IClientRegistry
	sentry     sentry.ISentry
	breaker    circuitbreaker.ICircuitBreaker
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
	breaker circuitbreaker.ICircuitBreaker,
) ISubOrderService {
	return &SubOrder{
		repository: repository,
		client:     client,
		sentry:     sentry,
		breaker:    breaker,
	}
}

func (o *SubOrder) GetSubOrderList(
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

	subOrders, total, err := o.repository.GetSubOrder().FindAllWithPagination(ctx, request)
	if err != nil {
		return nil, err
	}

	orderResponses := make([]subOrderDTO.SubOrderResponse, 0, len(orders))
	for _, subOrder := range subOrders {
		orderResponses = append(orderResponses, subOrderDTO.SubOrderResponse{
			OrderID:      subOrder.Order.UUID,
			SubOrderID:   subOrder.UUID,
			SubOrderName: subOrder.SubOrderName,
			CustomerID:   subOrder.Order.CustomerID,
			PackageID:    subOrder.Order.PackageID,
			Amount:       subOrder.Amount,
			Status:       subOrder.Status,
			IsPaid:       subOrder.IsPaid,
			CreatedAt:    subOrder.CreatedAt,
			UpdatedAt:    subOrder.UpdatedAt,
			Payment: &orderPaymentDTO.OrderPaymentResponse{
				PaymentID:   subOrder.Payment.PaymentID,
				PaymentLink: *subOrder.Payment.PaymentURL,
				Status:      subOrder.Payment.Status,
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

func (o *SubOrder) GetOrderDetail(ctx context.Context, subOrderUUID string) (*subOrderDTO.SubOrderResponse, error) {
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
		OrderID:      subOrder.Order.UUID,
		SubOrderID:   subOrder.UUID,
		SubOrderName: subOrder.SubOrderName,
		CustomerID:   subOrder.Order.CustomerID,
		PackageID:    subOrder.Order.PackageID,
		Amount:       subOrder.Amount,
		Status:       subOrder.Status,
		IsPaid:       subOrder.IsPaid,
		Payment: &orderPaymentDTO.OrderPaymentResponse{
			PaymentID:   subOrder.Payment.PaymentID,
			PaymentLink: *subOrder.Payment.PaymentURL,
			Status:      subOrder.Payment.Status,
		},
	}
	return response, nil
}

func (o *SubOrder) CreateOrder(
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

func (o *SubOrder) randomNumber() int {
	random := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
	number := random.Intn(1000000)
	return number
}

//nolint:cyclop,gocognit
func (o *SubOrder) createDownPaymentOrder(
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
		packageResponse  *packageClient.PackageData
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

	tx := o.repository.GetTx()
	err = tx.Transaction(func(tx *gorm.DB) error {
		packageRequest := circuitbreaker.BreakerFunc(func() (interface{}, error) {
			packageResponse, txErr = o.client.GetPackage().GetDetailPackage(ctx, request.PackageID.String())
			if txErr != nil {
				return nil, txErr
			}

			return packageResponse, nil
		})
		txErr = o.breaker.Execute(ctx, packageRequest)
		if txErr != nil {
			return txErr
		}

		total := float64(packageResponse.Price) * float64(packageResponse.MinimumDownPayment) / 100
		if total != request.Amount {
			newError := fmt.Errorf("down payment must be %d%% from package price", packageResponse.MinimumDownPayment) //nolint:goerr113,lll
			return newError
		}

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

		rbacRequest := circuitbreaker.BreakerFunc(func() (interface{}, error) {
			customerResponse, txErr = o.client.GetRBAC().GetUserRBAC(ctx, request.CustomerID.String())
			if txErr != nil {
				return nil, txErr
			}

			return customerResponse, nil
		})

		txErr = o.breaker.Execute(ctx, rbacRequest)
		if txErr != nil {
			return txErr
		}

		order, txErr = o.repository.GetOrder().Create(ctx, tx, &orderDTO.OrderRequest{
			CustomerID:                 request.CustomerID.String(),
			CustomerName:               customerResponse.Name,
			CustomerEmail:              customerResponse.Email,
			CustomerPhone:              customerResponse.PhoneNumber,
			PackageID:                  request.PackageID.String(),
			RemainingOutstandingAmount: float64(packageResponse.Price),
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

		expiredAt := time.Now().Add(24 * time.Hour)
		paymentRequest := circuitbreaker.BreakerFunc(func() (interface{}, error) {
			paymentResponse, txErr = o.generatePaymentLink(
				ctx,
				subOrder,
				order,
				request,
				expiredAt,
			)
			if txErr != nil {
				return nil, txErr
			}

			return paymentResponse, nil
		})
		txErr = o.breaker.Execute(ctx, paymentRequest)
		if txErr != nil {
			return txErr
		}

		txErr = o.repository.GetOrderPayment().
			Create(ctx, tx, &orderPaymentDTO.OrderPaymentRequest{
				Amount:      request.Amount,
				SubOrderID:  subOrder.ID,
				PaymentID:   paymentResponse.UUID,
				PaymentLink: paymentResponse.PaymentLink,
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
		OrderID:      order.UUID,
		SubOrderID:   subOrder.UUID,
		SubOrderName: subOrder.SubOrderName,
		CustomerID:   order.CustomerID,
		PackageID:    order.PackageID,
		Amount:       subOrder.Amount,
		Status:       subOrder.Status,
		OrderDate:    subOrder.OrderDate,
		IsPaid:       subOrder.IsPaid,
		Payment: &orderPaymentDTO.OrderPaymentResponse{
			PaymentID:   paymentResponse.UUID,
			PaymentLink: paymentResponse.PaymentLink,
			Status:      paymentResponse.Status,
		},
	}
	return &response, nil
}

//nolint:cyclop
func (o *SubOrder) createHalfPaymentOrder(
	ctx context.Context,
	request *subOrderDTO.SubOrderRequest,
) (*subOrderDTO.SubOrderResponse, error) {
	const logCtx = "services.suborder.sub_order.createHalfPaymentOrder"
	var (
		subOrder        *models.SubOrder
		order           *models.Order
		txErr           error
		paymentResponse *paymentClient.PaymentData
		err             error
		orderHistories  []orderHistoryDTO.OrderHistoryRequest
		span            = o.sentry.StartSpan(ctx, logCtx)
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

		expiredAt := time.Now().Add(24 * time.Hour)
		paymentRequest := circuitbreaker.BreakerFunc(func() (interface{}, error) {
			paymentResponse, txErr = o.generatePaymentLink(
				ctx,
				subOrder,
				order,
				request,
				expiredAt,
			)
			if txErr != nil {
				return nil, txErr
			}

			return paymentResponse, nil
		})
		txErr = o.breaker.Execute(ctx, paymentRequest)
		if txErr != nil {
			return txErr
		}

		txErr = o.repository.GetOrderPayment().
			Create(ctx, tx, &orderPaymentDTO.OrderPaymentRequest{
				Amount:      request.Amount,
				SubOrderID:  subOrder.ID,
				PaymentID:   paymentResponse.UUID,
				PaymentLink: paymentResponse.PaymentLink,
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
		OrderID:      order.UUID,
		SubOrderID:   subOrder.UUID,
		SubOrderName: subOrder.SubOrderName,
		CustomerID:   order.CustomerID,
		PackageID:    order.PackageID,
		Amount:       subOrder.Amount,
		Status:       subOrder.Status,
		OrderDate:    subOrder.OrderDate,
		IsPaid:       subOrder.IsPaid,
		Payment: &orderPaymentDTO.OrderPaymentResponse{
			PaymentID:   paymentResponse.UUID,
			PaymentLink: paymentResponse.PaymentLink,
			Status:      paymentResponse.Status,
		},
	}
	return &response, nil
}

//nolint:cyclop
func (o *SubOrder) createFullPaymentOrder(
	ctx context.Context,
	request *subOrderDTO.SubOrderRequest,
) (*subOrderDTO.SubOrderResponse, error) {
	const logCtx = "services.suborder.sub_order.createFullPaymentOrder"
	var (
		subOrder        *models.SubOrder
		order           *models.Order
		txErr           error
		paymentResponse *paymentClient.PaymentData
		err             error
		orderHistories  []orderHistoryDTO.OrderHistoryRequest
		span            = o.sentry.StartSpan(ctx, logCtx)
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

		expiredAt := time.Now().Add(24 * time.Hour)
		paymentRequest := circuitbreaker.BreakerFunc(func() (interface{}, error) {
			paymentResponse, txErr = o.generatePaymentLink(
				ctx,
				subOrder,
				order,
				request,
				expiredAt,
			)
			if txErr != nil {
				return nil, txErr
			}

			return paymentResponse, nil
		})
		txErr = o.breaker.Execute(ctx, paymentRequest)
		if txErr != nil {
			return txErr
		}

		txErr = o.repository.GetOrderPayment().
			Create(ctx, tx, &orderPaymentDTO.OrderPaymentRequest{
				Amount:      request.Amount,
				SubOrderID:  subOrder.ID,
				PaymentID:   paymentResponse.UUID,
				PaymentLink: paymentResponse.PaymentLink,
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
		OrderID:      order.UUID,
		SubOrderID:   subOrder.UUID,
		SubOrderName: subOrder.SubOrderName,
		CustomerID:   order.CustomerID,
		PackageID:    order.PackageID,
		Amount:       subOrder.Amount,
		Status:       subOrder.Status,
		OrderDate:    subOrder.OrderDate,
		IsPaid:       subOrder.IsPaid,
		Payment: &orderPaymentDTO.OrderPaymentResponse{
			PaymentID:   paymentResponse.UUID,
			PaymentLink: paymentResponse.PaymentLink,
			Status:      paymentResponse.Status,
		},
	}
	return &response, nil
}

func (o *SubOrder) generatePaymentLink(
	ctx context.Context,
	subOrder *models.SubOrder,
	order *models.Order,
	request *subOrderDTO.SubOrderRequest,
	expiredAt time.Time,
) (*paymentClient.PaymentData, error) {
	paymentResponse, err := o.client.GetPayment().CreatePaymentLink(ctx, &paymentClient.PaymentRequest{
		OrderID:     subOrder.UUID,
		ExpiredAt:   expiredAt,
		Amount:      request.Amount,
		Description: request.PaymentType.Title(),
		CustomerDetail: paymentClient.CustomerDetail{
			Name:  order.CustomerName,
			Email: order.CustomerEmail,
			Phone: order.CustomerPhone,
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
	if err != nil {
		return nil, err
	}

	return paymentResponse, nil
}

func (o *SubOrder) generateInvoice(
	ctx context.Context,
	request *invoiceModel.InvoiceRequest,
) (*invoiceModel.InvoiceData, error) {
	invoiceResponse, err := o.client.GetInvoice().GenerateInvoice(ctx, request)
	if err != nil {
		return nil, err
	}

	return invoiceResponse, nil
}

func (o *SubOrder) Cancel(ctx context.Context, subOrderUUID string) error {
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

//nolint:cyclop,gocognit
func (o *SubOrder) processPayment(
	ctx context.Context,
	request *subOrderDTO.PaymentRequest,
	status constant.OrderStatus,
) error {
	const logCtx = "services.suborder.sub_order.processPayment"
	var (
		updateRequest       subOrderDTO.UpdateSubOrderRequest
		paymentResult       *models.OrderPayment
		invoiceResponse     *invoiceModel.InvoiceData
		allSubOrder         []models.SubOrder
		paidAt, completedAt *time.Time
		isPaid              = false
		order               *models.Order
		total               float64
		indonesianTitle     string
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

			paymentResult, txErr = o.repository.GetOrderPayment().FindByPaymentID(ctx, tx, request.PaymentID.String())
			if txErr != nil {
				return txErr
			}

			allSubOrder, txErr = o.repository.GetSubOrder().FindAllByOrderID(ctx, order.ID)
			if txErr != nil {
				return txErr
			}
			items := make([]invoiceModel.Item, 0, len(allSubOrder))
			var totalPrice float64
			for _, item := range allSubOrder {
				switch item.PaymentType {
				case constant.PTDownPayment:
					indonesianTitle = constant.PTDownPaymentIndonesianTitle.String()
				case constant.PTHalfPayment:
					indonesianTitle = constant.PTHalfPaymentIndonesianTitle.String()
				case constant.PTFullPayment:
					indonesianTitle = constant.PTFullPaymentIndonesianTitle.String()
				}

				totalPrice += item.Amount
				items = append(items, invoiceModel.Item{
					Description: indonesianTitle,
					Price:       helper.RupiahFormat(&item.Amount),
				})
			}

			if total == 0 {
				isPaid = true
			} else {
				isPaid = false
			}
			invoiceNumber := fmt.Sprintf("INV/%s/ORD/%d", time.Now().Format("20060102"), o.randomNumber())
			invoiceRequest := circuitbreaker.BreakerFunc(func() (interface{}, error) {
				paidDay := paymentResult.PaidAt.Format("02")
				paidMonth := helper.ConvertToIndonesianMonth(paymentResult.PaidAt.Format("January"))
				paidYear := paymentResult.PaidAt.Format("2006")
				paymentMethod := helper.Ucwords(strings.ReplaceAll(*paymentResult.PaymentType, "_", " "))
				invoiceResponse, txErr = o.generateInvoice(
					ctx,
					&invoiceModel.InvoiceRequest{
						InvoiceNumber: invoiceNumber,
						TemplateID:    config.Config.InternalService.Invoice.TemplateID,
						CreatedBy:     order.CustomerID,
						Data: invoiceModel.Data{
							Customer: invoiceModel.Customer{
								Name:        order.CustomerName,
								Email:       order.CustomerEmail,
								PhoneNumber: order.CustomerPhone,
							},
							PaymentDetail: invoiceModel.PaymentDetail{
								PaymentMethod:              paymentMethod,
								BankName:                   strings.ToUpper(*paymentResult.Bank),
								VaNumber:                   *paymentResult.VANumber,
								RemainingOutstandingAmount: helper.RupiahFormat(&total),
								Date:                       fmt.Sprintf("%s %s %s", paidDay, paidMonth, paidYear),
								IsPaid:                     isPaid,
							},
							Items: items,
							Total: helper.RupiahFormat(&totalPrice),
						},
					},
				)
				if txErr != nil {
					return nil, txErr
				}

				return invoiceResponse, nil
			})
			txErr = o.breaker.Execute(ctx, invoiceRequest)
			if txErr != nil {
				return txErr
			}

			txErr = o.repository.GetOrderInvoice().Create(ctx, tx, &models.OrderInvoice{
				SubOrderID:    subOrder.ID,
				InvoiceID:     invoiceResponse.UUID,
				InvoiceNumber: invoiceNumber,
				InvoiceURL:    invoiceResponse.URL,
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

func (o *SubOrder) ReceivePendingPayment(ctx context.Context, request *subOrderDTO.PaymentRequest) error {
	return o.processPayment(ctx, request, constant.PendingPayment)
}

func (o *SubOrder) ReceivePaymentSettlement(ctx context.Context, request *subOrderDTO.PaymentRequest) error {
	return o.processPayment(ctx, request, constant.PaymentSuccess)
}

func (o *SubOrder) ReceivePaymentExpire(ctx context.Context, request *subOrderDTO.PaymentRequest) error {
	return o.processPayment(ctx, request, constant.Cancelled)
}
