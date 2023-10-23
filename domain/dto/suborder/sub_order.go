package dto

import (
	"github.com/google/uuid"

	"order-service/constant"
	orderPaymentDTO "order-service/domain/dto/orderpayment"

	"time"
)

type SubOrderRequest struct {
	OrderID     uint                 `json:"orderID"`
	CustomerID  string               `json:"customerID" validate:"required"`
	PackageID   string               `json:"packageID" validate:"required"`
	Amount      float64              `json:"amount" validate:"required"`
	OrderDate   time.Time            `json:"orderDate" validate:"required"`
	Status      constant.OrderStatus `json:"status"`
	IsPaid      *bool                `json:"isPaid"`
	PaymentType constant.PaymentType `json:"paymentType" validate:"required,oneof=down_payment half_payment full_payment"`
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
	UUID         uuid.UUID                             `json:"orderID"`
	SubOrderName string                                `json:"subOrderName"`
	CustomerID   string                                `json:"customerID"`
	PackageID    string                                `json:"packageID"`
	Amount       float64                               `json:"amount"`
	Status       constant.OrderStatus                  `json:"status"`
	OrderDate    time.Time                             `json:"orderDate,omitempty"`
	IsPaid       *bool                                 `json:"isPaid"`
	CompletedAt  *time.Time                            `json:"completedAt,omitempty"`
	CanceledAt   *time.Time                            `json:"canceledAt,omitempty"`
	CreatedAt    *time.Time                            `json:"createdAt"`
	UpdatedAt    *time.Time                            `json:"updatedAt"`
	Payment      *orderPaymentDTO.OrderPaymentResponse `json:"payment"`
}
