package dto

import (
	"github.com/google/uuid"

	"order-service/constant"
	orderPaymentDTO "order-service/domain/dto/orderpayment"

	"time"
)

type SubOrderRequest struct {
	OrderID     uuid.UUID            `json:"orderID" validate:"required_unless=PaymentType down_payment"`
	CustomerID  uuid.UUID            `json:"customerID" validate:"required"`
	PackageID   uuid.UUID            `json:"packageID" validate:"required"`
	Amount      float64              `json:"amount" validate:"required"`
	OrderDate   time.Time            `json:"orderDate" validate:"required"`
	Status      constant.OrderStatus `json:"status"`
	IsPaid      *bool                `json:"isPaid"`
	PaymentType constant.PaymentType `json:"paymentType" validate:"required,oneof=down_payment half_payment full_payment"`
	CanceledAt  *time.Time           `json:"canceledAt"`
}

type UpdateSubOrderRequest struct {
	OrderID    uuid.UUID            `json:"order_id"`
	Status     constant.OrderStatus `json:"status"`
	IsPaid     *bool                `json:"is_paid"`
	CanceledAt *time.Time           `json:"canceled_at"`
}

type PaymentRequest struct {
	OrderID     uuid.UUID  `json:"order_id"`
	PaymentID   uuid.UUID  `json:"payment_id"`
	PaymentLink string     `json:"payment_link"`
	PaymentType string     `json:"payment_type"`
	Amount      float64    `json:"amount"`
	Status      string     `json:"status"`
	VaNumber    *string    `json:"va_number"`
	Bank        *string    `json:"bank"`
	Acquirer    *string    `json:"acquirer"`
	ExpiredAt   *time.Time `json:"expired_at"`
	PaidAt      *time.Time `json:"paid_at"`
}

type CancelRequest struct {
	UUID   uuid.UUID            `json:"uuid,omitempty"`
	Status constant.OrderStatus `json:"status"`
}

type SubOrderRequestParam struct {
	Page  int `form:"page" validate:"required"`
	Limit int `form:"limit" validate:"required"`
}

type SubOrderResponse struct {
	OrderID      uuid.UUID                             `json:"orderID"`
	SubOrderID   uuid.UUID                             `json:"subOrderID"`
	SubOrderName string                                `json:"subOrderName"`
	CustomerID   string                                `json:"customerID"`
	PackageID    string                                `json:"packageID"`
	Amount       float64                               `json:"amount"`
	Status       constant.OrderStatus                  `json:"status"`
	OrderDate    time.Time                             `json:"orderDate,omitempty"`
	IsPaid       *bool                                 `json:"isPaid"`
	CanceledAt   *time.Time                            `json:"canceledAt,omitempty"`
	CreatedAt    *time.Time                            `json:"createdAt"`
	UpdatedAt    *time.Time                            `json:"updatedAt"`
	Payment      *orderPaymentDTO.OrderPaymentResponse `json:"payment"`
}
