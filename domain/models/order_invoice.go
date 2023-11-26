package models

import (
	"github.com/google/uuid"

	"time"
)

type OrderInvoice struct {
	ID            uint `gorm:"primaryKey;autoIncrement"`
	SubOrderID    uint
	InvoiceID     uuid.UUID
	InvoiceNumber string
	InvoiceURL    string
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}
