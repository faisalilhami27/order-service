package clients

import (
	clientConfig "order-service/clients/config"
	invoiceClient "order-service/clients/invoice"
	paymentClient "order-service/clients/payment"
	"order-service/common/sentry"
	"order-service/config"
)

type Client struct {
	sentry sentry.ISentry
}

type IClientRegistry interface {
	GetPayment() paymentClient.IPaymentClient
	GetInvoice() invoiceClient.IInvoiceClient
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

func (c *Client) GetInvoice() invoiceClient.IInvoiceClient {
	return invoiceClient.NewInvoiceClient(
		c.sentry,
		clientConfig.NewClientConfig(
			clientConfig.WithBaseURL(config.Config.InternalService.Invoice.Host),
			clientConfig.WithSecretKey(config.Config.InternalService.Invoice.SecretKey),
		))
}
