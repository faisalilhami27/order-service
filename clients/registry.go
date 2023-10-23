package clients

import paymentClient "order-service/clients/payment"

type Client struct{}

type IClientRegistry interface {
	GetPayment() paymentClient.IPaymentClient
}

func NewClientRegistry() IClientRegistry {
	return &Client{}
}

func (c *Client) GetPayment() paymentClient.IPaymentClient {
	return paymentClient.NewPaymentClient()
}
