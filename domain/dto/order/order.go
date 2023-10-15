package order

import (
	"github.com/google/uuid"

	"order-service/constant"
	orderPaymentDTO "order-service/domain/dto/orderpayment"
	orderModel "order-service/domain/models/order"

	"time"
)

type OrderRequest struct { //nolint:revive
	CustomerID string               `json:"customerID" validate:"required"`
	PackageID  string               `json:"packageID" validate:"required"`
	Amount     float64              `json:"amount" validate:"required"`
	OrderDate  time.Time            `json:"orderDate" validate:"required"`
	Status     constant.OrderStatus `json:"status"`
	IsPaid     *bool                `json:"isPaid"`
}

type CancelRequest struct {
	UUID   uuid.UUID            `json:"uuid,omitempty"`
	Status constant.OrderStatus `json:"status"`
}

type OrderRequestParam struct { //nolint:revive
	Page  int `form:"page" validate:"required"`
	Limit int `form:"limit" validate:"required"`
}

type OrderResponse struct { //nolint:revive
	UUID        uuid.UUID                             `json:"orderID"`
	OrderName   string                                `json:"orderName"`
	CustomerID  string                                `json:"customerID"`
	PackageID   string                                `json:"packageID"`
	Amount      float64                               `json:"amount"`
	Status      constant.OrderStatus                  `json:"status"`
	OrderDate   time.Time                             `json:"orderDate"`
	IsPaid      *bool                                 `json:"isPaid"`
	CompletedAt *time.Time                            `json:"completedAt"`
	CanceledAt  *time.Time                            `json:"canceledAt"`
	CreatedAt   *time.Time                            `json:"createdAt"`
	UpdatedAt   *time.Time                            `json:"updatedAt"`
	Payment     *orderPaymentDTO.OrderPaymentResponse `json:"payment"`
}

func ResponseFormatter(order *orderModel.Order, payment *orderPaymentDTO.OrderPaymentResponse) *OrderResponse {
	return &OrderResponse{
		UUID:        order.UUID,
		OrderName:   order.OrderName,
		CustomerID:  order.CustomerID,
		PackageID:   order.PackageID,
		Amount:      order.Amount,
		Status:      order.Status,
		OrderDate:   order.OrderDate,
		IsPaid:      order.IsPaid,
		CompletedAt: order.CompletedAt,
		CanceledAt:  order.CanceledAt,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
		Payment:     payment,
	}
}
