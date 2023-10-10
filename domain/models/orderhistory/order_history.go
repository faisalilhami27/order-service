package orderhistory

import (
	"order-service/constant"
	"time"
)

type OrderHistory struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	OrderID   uint
	Status    constant.OrderStatusString
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
