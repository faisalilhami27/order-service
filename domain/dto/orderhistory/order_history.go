package dto

import (
	"order-service/constant"
)

type OrderHistoryRequest struct {
	SubOrderID uint                       `json:"subOrderID"`
	Status     constant.OrderStatusString `json:"status"`
}
