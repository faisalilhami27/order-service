package clients

import (
	"github.com/google/uuid"

	"order-service/constant"

	"time"
)

type PaymentRequest struct {
	OrderID        uuid.UUID                 `json:"order_id"`
	ExpiredAt      time.Time                 `json:"expired_at"`
	Amount         float64                   `json:"amount"`
	Description    constant.PaymentTypeTitle `json:"description"`
	CustomerDetail CustomerDetail            `json:"customer_details"`
	ItemDetail     []ItemDetail              `json:"item_details"`
}

type CustomerDetail struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type ItemDetail struct {
	ID       uuid.UUID                 `json:"id"`
	Name     constant.PaymentTypeTitle `json:"name"`
	Amount   float64                   `json:"amount"`
	Quantity int                       `json:"quantity"`
}

type ErrorPaymentResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message any    `json:"message"`
}

type PaymentResponse struct {
	Code    int          `json:"code"`
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Data    PaymentData  `json:"data"`
	Error   *interface{} `json:"error,omitempty"`
}

type PaymentData struct {
	ID            int        `json:"id"`
	UUID          uuid.UUID  `json:"uuid"`
	OrderID       string     `json:"order_id"`
	Amount        float64    `json:"amount"`
	PaymentLink   string     `json:"payment_link"`
	Description   string     `json:"description"`
	ExpiredAt     time.Time  `json:"expired_at"`
	CreatedAt     *time.Time `json:"created_at"`
	Status        *string    `json:"status"`
	PaymentType   *string    `json:"payment_type"`
	VANumber      *string    `json:"va_number"`
	Bank          *string    `json:"bank"`
	Acquirer      *string    `json:"acquirer"`
	TransactionID *string    `json:"transaction_id"`
	BillerCode    *string    `json:"biller_code"`
	UpdatedAt     *time.Time `json:"updated_at"`
}
