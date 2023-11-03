package clients

import (
	clientConfig "order-service/clients/config"
	paymentClient "order-service/clients/payment"
	rbacClient "order-service/clients/rbac"
	"order-service/common/sentry"
	"order-service/config"
)

type Client struct {
	sentry sentry.ISentry
}

type IClientRegistry interface {
	GetPayment() paymentClient.IPaymentClient
	GetRBAC() rbacClient.IRbacClient
}

func NewClientRegistry(sentry sentry.ISentry) IClientRegistry {
	return &Client{
		sentry: sentry,
	}
}

func (c *Client) GetPayment() paymentClient.IPaymentClient {
	return paymentClient.NewPaymentClient(
		c.sentry,
		clientConfig.NewClientConfig(
			clientConfig.WithBaseURL(config.Config.InternalService.Payment.Host),
			clientConfig.WithSecretKey(config.Config.InternalService.Payment.SecretKey),
		))
}

func (c *Client) GetRBAC() rbacClient.IRbacClient {
	return rbacClient.NewRBACClient(
		c.sentry,
		clientConfig.NewClientConfig(
			clientConfig.WithBaseURL(config.Config.InternalService.RBAC.Host),
			clientConfig.WithSecretKey(config.Config.InternalService.RBAC.SecretKey),
		))
}
