package clients

type NotificationRequest struct {
	PhoneNumber string           `json:"phone_number"`
	TemplateID  string           `json:"template_id"`
	Data        SendWhatsappData `json:"data"`
}

type SendWhatsappData struct {
	OrderID     string `json:"order_id"`
	Description string `json:"description"`
	Amount      string `json:"amount"`
	PaymentLink string `json:"payment_link"`
	ExpiredAt   string `json:"expired_at"`
}

type ErrorNotificationResponse struct {
	Status  string `json:"status"`
	Message any    `json:"message"`
}

type NotificationResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
