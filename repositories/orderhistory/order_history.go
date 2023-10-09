package orderhistory

import (
	"context"
	"time"

	"gorm.io/gorm"

	orderHistoryDTO "order-service/domain/dto/orderhistory"
	orderHistoryModel "order-service/domain/models/orderhistory"
)

type IOrderHistory struct {
	db *gorm.DB
}

type IOrderHistoryRepository interface {
	Create(context.Context, *gorm.DB, *orderHistoryDTO.OrderHistoryRequest) error
}

func NewOrder(db *gorm.DB) IOrderHistoryRepository {
	return &IOrderHistory{
		db: db,
	}
}

func (o *IOrderHistory) Create(ctx context.Context, tx *gorm.DB, request *orderHistoryDTO.OrderHistoryRequest) error {
	location, _ := time.LoadLocation("Asia/Jakarta") //nolint:errcheck
	datetime := time.Now().In(location)

	orderHistory := orderHistoryModel.OrderHistory{
		OrderID:   request.OrderID,
		Status:    request.Status,
		CreatedAt: &datetime,
		UpdatedAt: &datetime,
	}
	err := tx.WithContext(ctx).Create(&orderHistory).Error
	if err != nil {
		return err
	}
	return nil
}
