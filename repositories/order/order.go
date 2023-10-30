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

	errorGeneral "order-service/constant/error"
	orderDTO "order-service/domain/dto/order"
	orderModel "order-service/domain/models/order"
	errorHelper "order-service/utils/error"
)

type IOrder struct {
	db *gorm.DB
}

type IOrderRepository interface {
	Create(context.Context, *gorm.DB, *orderDTO.OrderRequest) (*orderModel.Order, error)
	DeleteByOrderID(context.Context, *gorm.DB, uint) error
	FindOneOrderByUUID(context.Context, uuid.UUID) (*orderModel.Order, error)
	FindOneOrderByID(context.Context, uint) (*orderModel.Order, error)
	FindOneOrderByCustomerIDWithLocking(context.Context, uuid.UUID) (*orderModel.Order, error)
	Update(ctx context.Context, db *gorm.DB, request *orderDTO.OrderRequest) error
}

func NewOrder(db *gorm.DB) IOrderRepository {
	return &IOrder{
		db: db,
	}
}

func (o *IOrder) FindOneOrderByCustomerIDWithLocking(
	ctx context.Context,
	customerID uuid.UUID,
) (*orderModel.Order, error) {
	var order orderModel.Order
	err := o.db.WithContext(ctx).
		InnerJoins("INNER JOIN sub_orders ON sub_orders.order_id = orders.id").
		InnerJoins("INNER JOIN order_payments ON order_payments.sub_order_id = sub_orders.id").
		Where("order_payments.paid_at IS NULL AND sub_orders.canceled_at IS NULL").
		Where("customer_id = ?", customerID).
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

func (o *IOrder) FindOneOrderByUUID(
	ctx context.Context,
	uuid uuid.UUID,
) (*orderModel.Order, error) {
	var order orderModel.Order
	err := o.db.WithContext(ctx).
		Where("uuid = ?", uuid).
		Order("id DESC").
		First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errorHelper.WrapError(errorGeneral.ErrSQLError)
	}
	return &order, nil
}

func (o *IOrder) FindOneOrderByID(
	ctx context.Context,
	id uint,
) (*orderModel.Order, error) {
	var order orderModel.Order
	err := o.db.WithContext(ctx).
		Where("id = ?", id).
		Order("id DESC").
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
	location, _ := time.LoadLocation("Asia/Jakarta") //nolint:errcheck
	datetime := time.Now().In(location)
	orderName, err := o.autoNumber(ctx)
	if err != nil {
		return nil, err
	}

	order := orderModel.Order{
		UUID:                       uuid.New(),
		OrderName:                  *orderName,
		RemainingOutstandingAmount: request.RemainingOutstandingAmount,
		CustomerID:                 request.CustomerID,
		PackageID:                  request.PackageID,
		CreatedAt:                  &datetime,
		UpdatedAt:                  &datetime,
	}
	err = tx.WithContext(ctx).Create(&order).Error
	if err != nil {
		return nil, errorHelper.WrapError(errorGeneral.ErrSQLError)
	}
	return &order, nil
}

func (o *IOrder) DeleteByOrderID(ctx context.Context, tx *gorm.DB, orderID uint) error {
	err := tx.WithContext(ctx).
		Where("uuid = ?", orderID).
		Delete(&orderModel.Order{}).Error
	if err != nil {
		return errorHelper.WrapError(errorGeneral.ErrSQLError)
	}
	return nil
}

func (o *IOrder) Update(ctx context.Context, tx *gorm.DB, request *orderDTO.OrderRequest) error {
	err := tx.WithContext(ctx).
		Model(&orderModel.Order{}).
		Where("uuid = ?", request.OrderID).
		Updates(map[string]interface{}{
			"remaining_outstanding_amount": request.RemainingOutstandingAmount,
			"completed_at":                 request.CompletedAt,
		}).Error
	if err != nil {
		return errorHelper.WrapError(errorGeneral.ErrSQLError)
	}
	return nil
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
