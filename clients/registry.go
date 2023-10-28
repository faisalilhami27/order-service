package clients

import (
	clientConfig "order-service/clients/config"
	paymentClient "order-service/clients/payment"
	"order-service/config"
)

type Client struct{}

type IClientRegistry interface {
	GetPayment() paymentClient.IPaymentClient
}

func NewClientRegistry() IClientRegistry {
	return &Client{}
}

func (c *Client) GetPayment() paymentClient.IPaymentClient {
	return paymentClient.NewPaymentClient(clientConfig.NewClientConfig(
		clientConfig.WithBaseURL(config.Config.InternalService.Payment.Host),
		clientConfig.WithSecretKey(config.Config.InternalService.Payment.SecretKey),
	))
}
