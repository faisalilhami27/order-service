package models

import (
	"order-service/constant"

	"time"
)

type OrderHistory struct {
	ID         uint `gorm:"primaryKey;autoIncrement"`
	SubOrderID uint
	Status     constant.OrderStatusString
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}
