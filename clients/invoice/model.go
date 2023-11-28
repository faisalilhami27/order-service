package clients

import (
	"github.com/google/uuid"
)

type InvoiceRequest struct {
	InvoiceNumber string `json:"invoice_number"`
	TemplateID    string `json:"template_id"`
	Data          Data   `json:"data"`
	CreatedBy     string `json:"created_by"`
}

type Data struct {
	Customer      Customer      `json:"customer"`
	PaymentDetail PaymentDetail `json:"payment_detail"`
	Item          Item          `json:"item"`
}

type Customer struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

type PaymentDetail struct {
	BankName      string `json:"bank_name"`
	PaymentMethod string `json:"payment_method"`
	VaNumber      string `json:"va_number"`
	Date          string `json:"date"`
}

type Item struct {
	Description string `json:"description"`
	Price       string `json:"price"`
}

type ErrorInvoiceResponse struct {
	Status  string `json:"status"`
	Message any    `json:"message"`
}

type InvoiceResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    InvoiceData `json:"data"`
}

type InvoiceData struct {
	UUID uuid.UUID `json:"uuid"`
	URL  string    `json:"url"`
}
