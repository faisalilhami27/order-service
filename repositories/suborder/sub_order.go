package repositories

import (
	"context"
	"errors"
	"fmt"
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
	subOrderModel "order-service/domain/models/suborder"
	errorHelper "order-service/utils/error"
)

type ISubOrder struct {
	db *gorm.DB
}

type ISubOrderRepository interface {
	Create(context.Context, *gorm.DB, *subOrderDTO.SubOrderRequest) (*subOrderModel.SubOrder, error)
	FindOneSubOrderByCustomerIDWithLocking(context.Context, uuid.UUID) (*subOrderModel.SubOrder, error)
	FindOneByUUID(context.Context, string) (*subOrderModel.SubOrder, error)
	FindAllWithPagination(context.Context, *subOrderDTO.SubOrderRequestParam) ([]subOrderModel.SubOrder, int64, error)
	Cancel(context.Context, *gorm.DB, *subOrderDTO.CancelRequest, *subOrderModel.SubOrder) error
	BulkCreate(context.Context, *gorm.DB, []subOrderDTO.SubOrderRequest) ([]subOrderModel.SubOrder, error)
}

func NewSubOrder(db *gorm.DB) ISubOrderRepository {
	return &ISubOrder{
		db: db,
	}
}

func (o *ISubOrder) FindAllWithPagination(
	ctx context.Context,
	request *subOrderDTO.SubOrderRequestParam,
) ([]subOrderModel.SubOrder, int64, error) {
	var (
		order []subOrderModel.SubOrder
		total int64
	)
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
		return nil, 0, errorHelper.WrapError(errorGeneral.ErrSQLError)
	}

	err = o.db.WithContext(ctx).
		Model(&order).
		Count(&total).Error
	if err != nil {
		return nil, 0, errorHelper.WrapError(errorGeneral.ErrSQLError)
	}

	return order, total, nil
}

func (o *ISubOrder) FindOneByUUID(ctx context.Context, orderUUID string) (*subOrderModel.SubOrder, error) {
	var order subOrderModel.SubOrder
	err := o.db.WithContext(ctx).
		Preload("Payment").
		Where("uuid = ?", orderUUID).
		First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errOrder.ErrOrderNotFound
		}
		return nil, errorHelper.WrapError(errorGeneral.ErrSQLError)
	}
	return &order, nil
}

func (o *ISubOrder) FindOneSubOrderByCustomerIDWithLocking(
	ctx context.Context,
	customerID uuid.UUID,
) (*subOrderModel.SubOrder, error) {
	var order subOrderModel.SubOrder
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
		return nil, errorHelper.WrapError(errorGeneral.ErrSQLError)
	}
	return &order, nil
}

func (o *ISubOrder) Create(
	ctx context.Context,
	tx *gorm.DB,
	request *subOrderDTO.SubOrderRequest,
) (*subOrderModel.SubOrder, error) {
	isPaid := false
	subOrderName, err := o.autoNumber(ctx)
	if err != nil {
		return nil, err
	}

	subOrder := subOrderModel.SubOrder{
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

	err = tx.WithContext(ctx).Create(&subOrder).Error
	if err != nil {
		return nil, errorHelper.WrapError(errorGeneral.ErrSQLError)
	}
	return &subOrder, nil
}

func (o *ISubOrder) BulkCreate(
	ctx context.Context,
	tx *gorm.DB,
	requests []subOrderDTO.SubOrderRequest,
) ([]subOrderModel.SubOrder, error) {
	isPaid := false
	subOrderName, err := o.autoNumber(ctx)
	if err != nil {
		return nil, err
	}

	subOrders := make([]subOrderModel.SubOrder, 0, len(requests))
	for _, request := range requests {
		subOrder := subOrderModel.SubOrder{
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
		return nil, errorHelper.WrapError(errorGeneral.ErrSQLError)
	}

	return subOrders, nil
}

func (o *ISubOrder) Cancel(
	ctx context.Context,
	tx *gorm.DB,
	request *subOrderDTO.CancelRequest,
	current *subOrderModel.SubOrder,
) error {
	var (
		canceledAt = time.Now()
		order      subOrderModel.SubOrder
	)
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
		return errorHelper.WrapError(errorGeneral.ErrSQLError)
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
			return nil, errorHelper.WrapError(errorGeneral.ErrSQLError)
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
