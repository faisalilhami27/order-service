package order

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"order-service/constant"

	"time"
)

type Order struct {
	ID          int64                `db:"id" gorm:"primaryKey;autoIncrement"`
	UUID        uuid.UUID            `db:"uuid" gorm:"type:varchar(36);unique;not null"`
	OrderName   string               `db:"order_name" gorm:"type:varchar(20);unique;not null"`
	CustomerID  string               `db:"customer_id" gorm:"type:varchar(36);not null"`
	PackageID   string               `db:"package_id" gorm:"type:varchar(36);not null"`
	Amount      float64              `db:"amount" gorm:"not null"`
	Status      constant.OrderStatus `db:"status" gorm:"not null"`
	OrderDate   time.Time            `db:"order_date" gorm:"not null"`
	IsPaid      *bool                `db:"is_paid" gorm:"not null"`
	CompletedAt *time.Time           `db:"completed_at"`
	CanceledAt  *time.Time           `db:"canceled_at"`
	CreatedAt   *time.Time           `db:"created_at"`
	UpdatedAt   *time.Time           `db:"updated_at"`
	DeletedAt   *gorm.DeletedAt      `db:"deleted_at"`
}
