package clients

import (
	clientConfig "order-service/clients/config"
	invoiceClient "order-service/clients/invoice"
	notificationClient "order-service/clients/notification"
	paymentClient "order-service/clients/payment"
	weddingPackageClient "order-service/clients/weddingpackage"
	"order-service/common/sentry"
	"order-service/config"
)

type Client struct {
	sentry sentry.ISentry
}

type IClientRegistry interface {
	GetPayment() paymentClient.IPaymentClient
	GetWeddingPackage() weddingPackageClient.IWeddingPackageClient
	GetInvoice() invoiceClient.IInvoiceClient
	GetNotification() notificationClient.INotificationClient
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

func (c *Client) GetWeddingPackage() weddingPackageClient.IWeddingPackageClient {
	return weddingPackageClient.NewWeddingPackageClient(
		c.sentry,
		clientConfig.NewClientConfig(
			clientConfig.WithBaseURL(config.Config.InternalService.Package.Host),
			clientConfig.WithSecretKey(config.Config.InternalService.Package.SecretKey),
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

func (c *Client) GetNotification() notificationClient.INotificationClient {
	return notificationClient.NewNotificationClient(
		c.sentry,
		clientConfig.NewClientConfig(
			clientConfig.WithBaseURL(config.Config.InternalService.Notification.Host),
			clientConfig.WithSecretKey(config.Config.InternalService.Notification.SecretKey),
		))
}
