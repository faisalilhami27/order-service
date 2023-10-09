package orderhistory

import (
	"order-service/constant"
	"order-service/domain/models/order"
	"time"
)

type OrderHistory struct {
	ID        int64                      `db:"id" gorm:"primaryKey;autoIncrement"`
	OrderID   int64                      `db:"order_id"`
	Order     order.Order                `gorm:"foreignKey:OrderID"`
	Status    constant.OrderStatusString `db:"status"`
	CreatedAt *time.Time                 `db:"created_at"`
	UpdatedAt *time.Time                 `db:"updated_at"`
}
