package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"order-service/constant"
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
	OrderDate    time.Time            `gorm:"not null"`
	CanceledAt   *time.Time
	PaymentType  constant.PaymentType
	Order        Order          `gorm:"foreignKey:order_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Payment      OrderPayment   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Histories    []OrderHistory `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
	DeletedAt    *gorm.DeletedAt
}
