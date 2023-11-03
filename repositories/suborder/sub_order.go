package repositories

import (
	"context"
	"errors"
	"fmt"
	"order-service/common/sentry"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"order-service/common/state"
	"order-service/constant"
	errorGeneral "order-service/constant/error"
	errOrder "order-service/constant/error/order"
	subOrderDTO "order-service/domain/dto/suborder"
	subOrderModel "order-service/domain/models"
	errorHelper "order-service/utils/error"
)

type ISubOrder struct {
	db     *gorm.DB
	sentry sentry.ISentry
}

type ISubOrderRepository interface {
	Create(context.Context, *gorm.DB, *subOrderModel.SubOrder) (*subOrderModel.SubOrder, error)
	FindOneSubOrderByCustomerIDWithLocking(context.Context, uuid.UUID) (*subOrderModel.SubOrder, error)
	FindOneByUUID(context.Context, string) (*subOrderModel.SubOrder, error)
	FindOneByOrderIDAndPaymentType(context.Context, uint, string) (*subOrderModel.SubOrder, error)
	FindAllWithPagination(context.Context, *subOrderDTO.SubOrderRequestParam) ([]subOrderModel.SubOrder, int64, error)
	Cancel(context.Context, *gorm.DB, *subOrderDTO.CancelRequest, *subOrderModel.SubOrder) error
	BulkCreate(context.Context, *gorm.DB, []subOrderModel.SubOrder) ([]subOrderModel.SubOrder, error)
	Update(context.Context, *gorm.DB, *subOrderDTO.UpdateSubOrderRequest, *subOrderModel.SubOrder) error
}

func NewSubOrder(db *gorm.DB, sentry sentry.ISentry) ISubOrderRepository {
	return &ISubOrder{
		db:     db,
		sentry: sentry,
	}
}

func (o *ISubOrder) FindAllWithPagination(
	ctx context.Context,
	request *subOrderDTO.SubOrderRequestParam,
) ([]subOrderModel.SubOrder, int64, error) {
	const logCtx = "repositories.suborder.sub_order.FindAllWithPagination"
	var (
		span  = o.sentry.StartSpan(ctx, logCtx)
		order []subOrderModel.SubOrder
		total int64
	)
	ctx = o.sentry.SpanContext(span)
	defer o.sentry.Finish(span)

	limit := request.Limit
	offset := (request.Page - 1) * limit
	err := o.db.WithContext(ctx).
		Preload("Payment").
		Preload("Order").
		Limit(limit).
		Offset(offset).
		Find(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, errOrder.ErrOrderNotFound
		}
		return nil, 0, errorHelper.WrapError(errorGeneral.ErrSQLError, o.sentry)
	}

	err = o.db.WithContext(ctx).
		Model(&order).
		Count(&total).Error
	if err != nil {
		return nil, 0, errorHelper.WrapError(errorGeneral.ErrSQLError, o.sentry)
	}

	return order, total, nil
}

func (o *ISubOrder) FindOneByUUID(ctx context.Context, orderUUID string) (*subOrderModel.SubOrder, error) {
	const logCtx = "repositories.suborder.sub_order.FindOneByUUID"
	var (
		span  = o.sentry.StartSpan(ctx, logCtx)
		order subOrderModel.SubOrder
	)
	ctx = o.sentry.SpanContext(span)
	defer o.sentry.Finish(span)

	err := o.db.WithContext(ctx).
		Preload("Payment").
		Preload("Order").
		Where("uuid = ?", orderUUID).
		First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errOrder.ErrOrderNotFound
		}
		return nil, errorHelper.WrapError(errorGeneral.ErrSQLError, o.sentry)
	}
	return &order, nil
}

func (o *ISubOrder) FindOneByOrderIDAndPaymentType(
	ctx context.Context,
	orderID uint,
	paymentType string,
) (*subOrderModel.SubOrder, error) {
	const logCtx = "repositories.suborder.sub_order.FindOneByOrderIDAndPaymentType"
	var (
		span  = o.sentry.StartSpan(ctx, logCtx)
		order subOrderModel.SubOrder
	)
	ctx = o.sentry.SpanContext(span)
	defer o.sentry.Finish(span)

	err := o.db.WithContext(ctx).
		InnerJoins("INNER JOIN order_payments op ON op.sub_order_id = sub_orders.id").
		Where("sub_orders.payment_type = ?", paymentType).
		Where("order_id = ?", orderID).
		First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errorHelper.WrapError(errorGeneral.ErrSQLError, o.sentry)
	}
	return &order, nil
}

func (o *ISubOrder) FindOneSubOrderByCustomerIDWithLocking(
	ctx context.Context,
	customerID uuid.UUID,
) (*subOrderModel.SubOrder, error) {
	const logCtx = "repositories.suborder.sub_order.FindOneSubOrderByCustomerIDWithLocking"
	var (
		span  = o.sentry.StartSpan(ctx, logCtx)
		order subOrderModel.SubOrder
	)
	ctx = o.sentry.SpanContext(span)
	defer o.sentry.Finish(span)

	err := o.db.WithContext(ctx).
		Where("customer_id = ?", customerID).
		Where("completed_at IS NULL").
		Where("canceled_at IS NULL").
		Order("id DESC").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errorHelper.WrapError(errorGeneral.ErrSQLError, o.sentry)
	}
	return &order, nil
}

func (o *ISubOrder) Create(
	ctx context.Context,
	tx *gorm.DB,
	request *subOrderModel.SubOrder,
) (*subOrderModel.SubOrder, error) {
	const logCtx = "repositories.suborder.sub_order.Create"
	var (
		span     = o.sentry.StartSpan(ctx, logCtx)
		subOrder subOrderModel.SubOrder
	)
	ctx = o.sentry.SpanContext(span)
	defer o.sentry.Finish(span)

	isPaid := false
	subOrderName, err := o.autoNumber(ctx)
	if err != nil {
		return nil, err
	}

	subOrder = subOrderModel.SubOrder{
		UUID:         uuid.New(),
		SubOrderName: *subOrderName,
		OrderID:      request.OrderID,
		Status:       request.Status,
		Amount:       request.Amount,
		PaymentType:  request.PaymentType,
		OrderDate:    request.OrderDate,
		IsPaid:       &isPaid,
	}

	st := state.NewStatusState(constant.Initial)
	if st.FSM.Cannot(request.Status.String()) {
		errorStatus := fmt.Errorf("%w from %s to %s",
			errorGeneral.ErrInvalidStatusTransition,
			st.FSM.Current(),
			request.Status.String())
		return nil, errorStatus
	}

	err = tx.WithContext(ctx).Create(&subOrder).Error
	if err != nil {
		return nil, errorHelper.WrapError(errorGeneral.ErrSQLError, o.sentry)
	}
	return &subOrder, nil
}

func (o *ISubOrder) BulkCreate(
	ctx context.Context,
	tx *gorm.DB,
	requests []subOrderModel.SubOrder,
) ([]subOrderModel.SubOrder, error) {
	const logCtx = "repositories.suborder.sub_order.BulkCreate"
	var (
		span     = o.sentry.StartSpan(ctx, logCtx)
		subOrder subOrderModel.SubOrder
	)
	ctx = o.sentry.SpanContext(span)
	defer o.sentry.Finish(span)

	isPaid := false
	subOrderName, err := o.autoNumber(ctx)
	if err != nil {
		return nil, err
	}

	subOrders := make([]subOrderModel.SubOrder, 0, len(requests))
	for _, request := range requests {
		subOrder = subOrderModel.SubOrder{
			UUID:         uuid.New(),
			SubOrderName: *subOrderName,
			OrderID:      request.OrderID,
			Status:       request.Status,
			Amount:       request.Amount,
			PaymentType:  request.PaymentType,
			IsPaid:       &isPaid,
		}

		st := state.NewStatusState(constant.Initial)
		if st.FSM.Cannot(request.Status.String()) {
			errorStatus := fmt.Errorf("%w from %s to %s",
				errorGeneral.ErrInvalidStatusTransition,
				st.FSM.Current(),
				request.Status.String())
			return nil, errorStatus
		}

		subOrders = append(subOrders, subOrder)
	}

	err = tx.WithContext(ctx).Create(&subOrders).Error
	if err != nil {
		return nil, errorHelper.WrapError(errorGeneral.ErrSQLError, o.sentry)
	}

	return subOrders, nil
}

func (o *ISubOrder) Cancel(
	ctx context.Context,
	tx *gorm.DB,
	request *subOrderDTO.CancelRequest,
	current *subOrderModel.SubOrder,
) error {
	const logCtx = "repositories.suborder.sub_order.Cancel"
	var (
		span       = o.sentry.StartSpan(ctx, logCtx)
		canceledAt = time.Now()
		order      subOrderModel.SubOrder
	)
	ctx = o.sentry.SpanContext(span)
	defer o.sentry.Finish(span)

	order = subOrderModel.SubOrder{
		CanceledAt: &canceledAt,
	}

	st := state.NewStatusState(current.Status)
	if st.FSM.Cannot(request.Status.String()) {
		errorStatus := fmt.Errorf("%w from %s to %s",
			errorGeneral.ErrInvalidStatusTransition,
			st.FSM.Current(),
			request.Status.String())
		return errorStatus
	}

	err := tx.WithContext(ctx).
		Model(&order).
		Where("uuid = ?", request.UUID).
		Updates(subOrderModel.SubOrder{
			Status:     constant.Cancelled,
			CanceledAt: &canceledAt,
		}).Error
	if err != nil {
		return errorHelper.WrapError(errorGeneral.ErrSQLError, o.sentry)
	}
	return nil
}

func (o *ISubOrder) autoNumber(ctx context.Context) (*string, error) {
	var (
		subOrder *subOrderModel.SubOrder
		result   string
		today    = time.Now().Format("20060102")
	)
	err := o.db.WithContext(ctx).Order("id desc").First(&subOrder).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorHelper.WrapError(errorGeneral.ErrSQLError, o.sentry)
		}
	}

	if subOrder.ID != 0 {
		subOrderName := subOrder.SubOrderName
		splitOrderName, _ := strconv.Atoi(subOrderName[8:13]) //nolint:errcheck
		code := splitOrderName + 1
		result = fmt.Sprintf("SUB-ORD-%05d-%s", code, today)
	} else {
		result = fmt.Sprintf("SUB-ORD-%05d-%s", 1, today)
	}

	return &result, nil
}

func (o *ISubOrder) Update(
	ctx context.Context,
	tx *gorm.DB,
	request *subOrderDTO.UpdateSubOrderRequest,
	current *subOrderModel.SubOrder,
) error {
	const logCtx = "repositories.suborder.sub_order.Update"
	var (
		span     = o.sentry.StartSpan(ctx, logCtx)
		subOrder subOrderModel.SubOrder
	)
	ctx = o.sentry.SpanContext(span)
	defer o.sentry.Finish(span)

	subOrder = subOrderModel.SubOrder{
		Status:     request.Status,
		IsPaid:     request.IsPaid,
		CanceledAt: request.CanceledAt,
	}

	st := state.NewStatusState(current.Status)
	if st.FSM.Cannot(request.Status.String()) {
		errorStatus := fmt.Errorf("%w from %s to %s",
			errorGeneral.ErrInvalidStatusTransition,
			st.FSM.Current(),
			request.Status.String())
		return errorStatus
	}

	err := tx.WithContext(ctx).
		Where("uuid = ?", current.UUID).
		Updates(&subOrder).Error
	if err != nil {
		return errorHelper.WrapError(errorGeneral.ErrSQLError, o.sentry)
	}
	return nil
}
