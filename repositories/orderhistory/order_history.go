package repositories

import (
	"context"
	"time"

	"gorm.io/gorm"

	errorGeneral "order-service/constant/error"
	orderHistoryDTO "order-service/domain/dto/orderhistory"
	orderHistoryModel "order-service/domain/models/orderhistory"
	errorHelper "order-service/utils/error"
)

type IOrderHistory struct {
	db *gorm.DB
}

type IOrderHistoryRepository interface {
	Create(context.Context, *gorm.DB, *orderHistoryDTO.OrderHistoryRequest) error
	BulkCreate(context.Context, *gorm.DB, []orderHistoryDTO.OrderHistoryRequest) error
}

func NewOrderHistory(db *gorm.DB) IOrderHistoryRepository {
	return &IOrderHistory{
		db: db,
	}
}

func (o *IOrderHistory) Create(ctx context.Context, tx *gorm.DB, request *orderHistoryDTO.OrderHistoryRequest) error {
	location, _ := time.LoadLocation("Asia/Jakarta") //nolint:errcheck
	datetime := time.Now().In(location)

	orderHistory := orderHistoryModel.OrderHistory{
		SubOrderID: request.SubOrderID,
		Status:     request.Status,
		CreatedAt:  &datetime,
		UpdatedAt:  &datetime,
	}
	err := tx.WithContext(ctx).Create(&orderHistory).Error
	if err != nil {
		return errorHelper.WrapError(errorGeneral.ErrSQLError)
	}
	return nil
}

func (o *IOrderHistory) BulkCreate(
	ctx context.Context,
	tx *gorm.DB,
	requests []orderHistoryDTO.OrderHistoryRequest,
) error {
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
		return errorHelper.WrapError(errorGeneral.ErrSQLError)
	}

	return nil
}
