package models

import (
	"github.com/google/uuid"

	"time"
)

type OrderPayment struct {
	ID          uint `gorm:"primaryKey;autoIncrement"`
	Amount      float64
	SubOrderID  uint
	PaymentID   uuid.UUID
	PaymentURL  *string
	Status      *string
	PaidAt      *time.Time
	ExpiredAt   *time.Time
	PaymentType *string `gorm:"null"`
	VANumber    *string `gorm:"null"`
	Bank        *string `gorm:"null"`
	Acquirer    *string `gorm:"null"`
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}
