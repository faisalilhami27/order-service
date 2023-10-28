package dto

import (
	dto "order-service/domain/dto/kafka"
	"time"
)

type PaymentData struct {
	OrderID     string     `json:"order_id"`
	PaymentID   string     `json:"payment_id"`
	Amount      float64    `json:"amount"`
	PaymentLink string     `json:"payment_link"`
	PaymentType string     `json:"payment_type"`
	VANumber    *string    `json:"va_number"`
	Bank        *string    `json:"bank"`
	Acquirer    *string    `json:"acquirer"`
	Description *string    `json:"description"`
	Status      string     `json:"status"`
	ExpiredAt   string     `json:"expired_at"`
	PaidAt      *time.Time `json:"paid_at"`
}

type PaymentContent struct {
	Event dto.KafkaMessageEvent             `json:"event"`
	Meta  dto.KafkaMessageMeta              `json:"meta"`
	Body  dto.KafkaMessageBody[PaymentData] `json:"body"`
}
