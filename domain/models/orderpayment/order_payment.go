package models

import (
	"github.com/google/uuid"

	"time"
)

type OrderPayment struct {
	ID         uint `gorm:"primaryKey;autoIncrement"`
	SubOrderID uint
	PaymentID  uuid.UUID
	InvoiceID  uuid.UUID
	PaymentURL string
	Status     *string
	PaidAt     *time.Time
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}
