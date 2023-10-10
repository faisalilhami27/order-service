package orderpayment

import (
	"github.com/google/uuid"
)

type OrderPaymentRequest struct { //nolint:revive
	OrderID    uint      `json:"order_id"`
	PaymentID  uuid.UUID `json:"payment_id"`
	InvoiceID  uuid.UUID `json:"invoice_id"`
	PaymentURL string    `json:"payment_url"`
	Status     *string   `json:"status"`
}

type OrderPaymentResponse struct { //nolint:revive
	PaymentID  uuid.UUID `json:"paymentID"`
	InvoiceID  uuid.UUID `json:"invoiceID"`
	PaymentURL string    `json:"paymentURL"`
	Status     *string   `json:"status"`
}
