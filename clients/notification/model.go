package clients

type NotificationRequest struct {
	PhoneNumber string            `json:"phone_number"`
	TemplateID  string            `json:"template_id"`
	Title       *Title            `json:"title,omitempty"`
	Data        *SendWhatsappData `json:"data,omitempty"`
	Button      *Button           `json:"button,omitempty"`
	Footer      *string           `json:"footer,omitempty"`
}

type SendWhatsappData struct {
	OrderID     string `json:"order_id"`
	Description string `json:"description"`
	Amount      string `json:"amount"`
	ExpiredAt   string `json:"expired_at"`
}

type Button struct {
	URL  *URL  `json:"url,omitempty"`
	Call *Call `json:"call,omitempty"`
}

type URL struct {
	Display string `json:"display"`
	Link    string `json:"link"`
}

type Call struct {
	Display string `json:"display"`
	Phone   string `json:"phone"`
}

type Title struct {
	Type    string `json:"type"`
	Content string `json:"content"`
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
