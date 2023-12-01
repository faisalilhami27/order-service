package clients

import (
	"github.com/google/uuid"

	"time"
)

type PackageResponse[T any] struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

type ResponseData struct {
	PackageResponse[PackageData]
}

type Promo struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Discount  string `json:"discount"`
}

type PackagePromo struct {
	PackageID int   `json:"package_id"`
	PromoID   int   `json:"promo_id"`
	Promo     Promo `json:"promo"`
}

type PackageData struct {
	ID                 int          `json:"id"`
	UUID               uuid.UUID    `json:"uuid"`
	Name               string       `json:"name"`
	Description        string       `json:"description"`
	Price              int          `json:"price"`
	Pack               int          `json:"pack"`
	MinimalDownPayment int          `json:"minimalDownPayment"`
	IsActive           bool         `json:"isActive"`
	CreatedAt          time.Time    `json:"createdAt"`
	UpdatedAt          time.Time    `json:"updatedAt"`
	PackagePromo       PackagePromo `json:"packagePromo"`
}
