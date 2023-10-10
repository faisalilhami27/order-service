package orderhistory

import (
	"order-service/constant"
)

type OrderHistoryRequest struct { //nolint:revive
	OrderID uint                       `json:"order_id"`
	Status  constant.OrderStatusString `json:"status"`
}
