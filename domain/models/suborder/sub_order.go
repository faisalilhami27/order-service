package models

import (
	"github.com/google/uuid"

	"gorm.io/gorm"

	"order-service/constant"
	orderModel "order-service/domain/models/order"
	orderHistoryModel "order-service/domain/models/orderhistory"
	orderPaymentModel "order-service/domain/models/orderpayment"

	"time"
)

type SubOrder struct {
	ID           uint                 `gorm:"primaryKey;autoIncrement"`
	UUID         uuid.UUID            `gorm:"type:varchar(36);unique;not null"`
	OrderID      uint                 `gorm:"not null"`
	SubOrderName string               `gorm:"type:varchar(25);unique;not null"`
	Amount       float64              `gorm:"not null"`
	Status       constant.OrderStatus `gorm:"not null"`
	IsPaid       *bool                `gorm:"not null"`
	CompletedAt  *time.Time
	CanceledAt   *time.Time
	PaymentType  constant.PaymentType
	Order        orderModel.Order                 `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Payment      orderPaymentModel.OrderPayment   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Histories    []orderHistoryModel.OrderHistory `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
	DeletedAt    *gorm.DeletedAt
}
