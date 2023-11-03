package repositories

import (
	"context"
	"order-service/common/sentry"
	orderHistoryModel "order-service/domain/models"
	"time"

	"gorm.io/gorm"

	errorGeneral "order-service/constant/error"
	orderHistoryDTO "order-service/domain/dto/orderhistory"
	errorHelper "order-service/utils/error"
)

type IOrderHistory struct {
	db     *gorm.DB
	sentry sentry.ISentry
}

type IOrderHistoryRepository interface {
	Create(context.Context, *gorm.DB, *orderHistoryDTO.OrderHistoryRequest) error
	BulkCreate(context.Context, *gorm.DB, []orderHistoryDTO.OrderHistoryRequest) error
}

func NewOrderHistory(db *gorm.DB, sentry sentry.ISentry) IOrderHistoryRepository {
	return &IOrderHistory{
		db:     db,
		sentry: sentry,
	}
}

func (o *IOrderHistory) Create(ctx context.Context, tx *gorm.DB, request *orderHistoryDTO.OrderHistoryRequest) error {
	const logCtx = "repositories.orderhistory.order_history.Create"
	var (
		span         = o.sentry.StartSpan(ctx, logCtx)
		orderHistory orderHistoryModel.OrderHistory
	)
	ctx = o.sentry.SpanContext(span)
	defer o.sentry.Finish(span)

	location, _ := time.LoadLocation("Asia/Jakarta") //nolint:errcheck
	datetime := time.Now().In(location)

	orderHistory = orderHistoryModel.OrderHistory{
		SubOrderID: request.SubOrderID,
		Status:     request.Status,
		CreatedAt:  &datetime,
		UpdatedAt:  &datetime,
	}
	err := tx.WithContext(ctx).Create(&orderHistory).Error
	if err != nil {
		return errorHelper.WrapError(errorGeneral.ErrSQLError, o.sentry)
	}
	return nil
}

func (o *IOrderHistory) BulkCreate(
	ctx context.Context,
	tx *gorm.DB,
	requests []orderHistoryDTO.OrderHistoryRequest,
) error {
	const logCtx = "repositories.orderhistory.order_history.BulkCreate"
	var (
		span = o.sentry.StartSpan(ctx, logCtx)
	)
	ctx = o.sentry.SpanContext(span)
	defer o.sentry.Finish(span)

	location, _ := time.LoadLocation("Asia/Jakarta") //nolint:errcheck
	datetime := time.Now().In(location)

	orderHistoryRequest := make([]orderHistoryModel.OrderHistory, 0, len(requests))
	for _, request := range requests {
		orderHistory := orderHistoryModel.OrderHistory{
			SubOrderID: request.SubOrderID,
			Status:     request.Status,
			CreatedAt:  &datetime,
			UpdatedAt:  &datetime,
		}
		orderHistoryRequest = append(orderHistoryRequest, orderHistory)
	}

	err := tx.WithContext(ctx).Create(&orderHistoryRequest).Error
	if err != nil {
		return errorHelper.WrapError(errorGeneral.ErrSQLError, o.sentry)
	}

	return nil
}
