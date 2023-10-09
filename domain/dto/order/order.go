package order

import (
	"github.com/google/uuid"

	"order-service/constant"
	orderModel "order-service/domain/models/order"

	"time"
)

type OrderRequest struct { //nolint:revive
	CustomerID string               `json:"customer_id"`
	PackageID  string               `json:"package_id"`
	Amount     float64              `json:"amount"`
	OrderDate  time.Time            `json:"order_date"`
	Status     constant.OrderStatus `json:"status"`
	IsPaid     *bool                `json:"is_paid"`
}

type OrderResponse struct { //nolint:revive
	UUID        uuid.UUID            `json:"orderID"`
	OrderName   string               `json:"orderName"`
	CustomerID  string               `json:"customerID"`
	PackageID   string               `json:"packageID"`
	Amount      float64              `json:"amount"`
	Status      constant.OrderStatus `json:"status"`
	OrderDate   time.Time            `json:"orderDate"`
	IsPaid      *bool                `json:"IsPaid"`
	CompletedAt *time.Time           `json:"completedAt"`
	CanceledAt  *time.Time           `json:"canceledAt"`
	CreatedAt   *time.Time           `json:"createdAt"`
	UpdatedAt   *time.Time           `json:"updatedAt"`
}

func ResponseFormatter(order *orderModel.Order) *OrderResponse {
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
	}
}
