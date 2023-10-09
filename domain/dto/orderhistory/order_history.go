package orderhistory

import (
	"order-service/constant"
)

type OrderHistoryRequest struct { //nolint:revive
	OrderID int64                      `json:"order_id"`
	Status  constant.OrderStatusString `json:"status"`
}
