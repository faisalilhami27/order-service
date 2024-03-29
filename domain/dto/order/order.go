package dto

import "time"

type OrderRequest struct {
	CustomerID                 string     `json:"customerID" validate:"required"`
	CustomerName               string     `json:"customerName"`
	CustomerEmail              string     `json:"customerEmail"`
	CustomerPhone              string     `json:"customerPhone"`
	PackageID                  string     `json:"packageID" validate:"required"`
	RemainingOutstandingAmount float64    `json:"remainingOutstandingAmount"`
	OrderID                    string     `json:"orderID"`
	CompletedAt                *time.Time `json:"completedAt"`
	IsPaid                     *bool      `json:"isPaid"`
}
