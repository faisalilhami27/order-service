package order

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
	orderDTO "order-service/domain/dto/order"
	orderModel "order-service/domain/models/order"
	errorHelper "order-service/utils/error"
)

type IOrder struct {
	db *gorm.DB
}

type IOrderRepository interface {
	Create(context.Context, *gorm.DB, *orderDTO.OrderRequest) (*orderModel.Order, error)
	FindOneOrderByCustomerIDWithLocking(context.Context, uuid.UUID) (*orderModel.Order, error)
	FindOneByUUID(context.Context, string) (*orderModel.Order, error)
	FindAllWithPagination(context.Context, *orderDTO.OrderRequestParam) ([]orderModel.Order, int64, error)
}

func NewOrder(db *gorm.DB) IOrderRepository {
	return &IOrder{
		db: db,
	}
}

func (o *IOrder) FindAllWithPagination(
	ctx context.Context,
	request *orderDTO.OrderRequestParam,
) ([]orderModel.Order, int64, error) {
	var (
		order []orderModel.Order
		total int64
	)
	limit := request.Limit
	offset := (request.Page - 1) * limit
	err := o.db.WithContext(ctx).
		Preload("Payment").
		Limit(limit).
		Offset(offset).
		First(&order).Error
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

func (o *IOrder) FindOneByUUID(ctx context.Context, orderUUID string) (*orderModel.Order, error) {
	var order orderModel.Order
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

func (o *IOrder) FindOneOrderByCustomerIDWithLocking(
	ctx context.Context,
	customerID uuid.UUID,
) (*orderModel.Order, error) {
	var order orderModel.Order
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

func (o *IOrder) Create(ctx context.Context, tx *gorm.DB, request *orderDTO.OrderRequest) (*orderModel.Order, error) {
	isPaid := false
	orderName, err := o.autoNumber(ctx)
	if err != nil {
		return nil, err
	}

	order := orderModel.Order{
		UUID:       uuid.New(),
		OrderName:  *orderName,
		CustomerID: request.CustomerID,
		PackageID:  request.PackageID,
		Amount:     request.Amount,
		OrderDate:  request.OrderDate,
		Status:     request.Status,
		IsPaid:     &isPaid,
	}

	st := state.NewStatusState(constant.Initial)
	if st.FSM.Cannot(request.Status.String()) {
		errorStatus := fmt.Errorf("%w from %s to %s",
			errorGeneral.ErrInvalidStatusTransition,
			st.FSM.Current(),
			request.Status.String())
		return nil, errorStatus
	}

	err = tx.WithContext(ctx).Create(&order).Error
	if err != nil {
		return nil, errorHelper.WrapError(errorGeneral.ErrSQLError)
	}
	return &order, nil
}

func (o *IOrder) autoNumber(ctx context.Context) (*string, error) {
	var (
		order  *orderModel.Order
		result string
		today  = time.Now().Format("20060102")
	)
	err := o.db.WithContext(ctx).Order("id desc").First(&order).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorHelper.WrapError(errorGeneral.ErrSQLError)
		}
	}

	if order.ID != 0 {
		orderName := order.OrderName
		splitOrderName, _ := strconv.Atoi(orderName[4:9]) //nolint:errcheck
		code := splitOrderName + 1
		result = fmt.Sprintf("ORD-%05d-%s", code, today)
	} else {
		result = fmt.Sprintf("ORD-%05d-%s", 1, today)
	}

	return &result, nil
}
