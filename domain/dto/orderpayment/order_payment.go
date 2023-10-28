package dto

import (
	"github.com/google/uuid"

	"time"
)

type OrderPaymentRequest struct {
	Amount      float64    `json:"amount"`
	SubOrderID  uint       `json:"sub_order_id"`
	PaymentID   uuid.UUID  `json:"payment_id"`
	InvoiceID   uuid.UUID  `json:"invoice_id,omitempty"`
	PaymentLink string     `json:"payment_link"`
	Status      *string    `json:"status"`
	PaymentType *string    `json:"payment_type"`
	VANumber    *string    `json:"va_number,omitempty"`
	Bank        *string    `json:"bank,omitempty"`
	Acquirer    *string    `json:"acquirer,omitempty"`
	ExpiredAt   *time.Time `json:"expired_at,omitempty"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
}

type OrderPaymentResponse struct {
	PaymentID   uuid.UUID `json:"paymentID"`
	InvoiceID   uuid.UUID `json:"invoiceID,omitempty"`
	PaymentLink string    `json:"paymentLink"`
	Status      *string   `json:"status"`
}
