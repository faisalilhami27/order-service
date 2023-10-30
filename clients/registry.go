package clients

import (
	clientConfig "order-service/clients/config"
	paymentClient "order-service/clients/payment"
	rbacClient "order-service/clients/rbac"
	"order-service/config"
)

type Client struct{}

type IClientRegistry interface {
	GetPayment() paymentClient.IPaymentClient
	GetRBAC() rbacClient.IRbacClient
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

func (c *Client) GetRBAC() rbacClient.IRbacClient {
	return rbacClient.NewRBACClient(clientConfig.NewClientConfig(
		clientConfig.WithBaseURL(config.Config.InternalService.RBAC.Host),
		clientConfig.WithSecretKey(config.Config.InternalService.RBAC.SecretKey),
	))
}
