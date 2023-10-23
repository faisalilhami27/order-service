package models

import (
	"github.com/google/uuid"

	"gorm.io/gorm"

	"time"
)

type Order struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	UUID        uuid.UUID `gorm:"type:varchar(36);unique;not null"`
	OrderName   string    `gorm:"type:varchar(20);unique;not null"`
	CustomerID  string    `gorm:"type:varchar(36);not null"`
	PackageID   string    `gorm:"type:varchar(36);not null"`
	OrderDate   time.Time `gorm:"not null"`
	CompletedAt *time.Time
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	DeletedAt   *gorm.DeletedAt
}
