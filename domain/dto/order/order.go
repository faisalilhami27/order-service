package dto

import "time"

type OrderRequest struct {
	CustomerID  string     `json:"customerID" validate:"required"`
	PackageID   string     `json:"packageID" validate:"required"`
	OrderDate   time.Time  `json:"orderDate" validate:"required"`
	CompletedAt *time.Time `json:"completedAt"`
}
