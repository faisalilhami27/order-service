package order

import (
	"github.com/google/uuid"

	"gorm.io/gorm"

	"order-service/constant"
	orderHistoryModel "order-service/domain/models/orderhistory"
	orderPaymentModel "order-service/domain/models/orderpayment"

	"time"
)

type Order struct {
	ID          uint                 `gorm:"primaryKey;autoIncrement"`
	UUID        uuid.UUID            `gorm:"type:varchar(36);unique;not null"`
	OrderName   string               `gorm:"type:varchar(20);unique;not null"`
	CustomerID  string               `gorm:"type:varchar(36);not null"`
	PackageID   string               `gorm:"type:varchar(36);not null"`
	Amount      float64              `gorm:"not null"`
	Status      constant.OrderStatus `gorm:"not null"`
	OrderDate   time.Time            `gorm:"not null"`
	IsPaid      *bool                `gorm:"not null"`
	CompletedAt *time.Time
	CanceledAt  *time.Time
	Payment     orderPaymentModel.OrderPayment   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Histories   []orderHistoryModel.OrderHistory `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	DeletedAt   *gorm.DeletedAt
}
