package dto

import (
	"github.com/google/uuid"
)

type OrderPaymentRequest struct {
	SubOrderID  uint      `json:"subOrderID"`
	PaymentID   uuid.UUID `json:"paymentID"`
	InvoiceID   uuid.UUID `json:"invoiceID,omitempty"`
	PaymentLink string    `json:"paymentLink"`
	Status      *string   `json:"status"`
}

type OrderPaymentResponse struct {
	PaymentID   uuid.UUID `json:"paymentID"`
	InvoiceID   uuid.UUID `json:"invoiceID,omitempty"`
	PaymentLink string    `json:"paymentLink"`
	Status      *string   `json:"status"`
}
