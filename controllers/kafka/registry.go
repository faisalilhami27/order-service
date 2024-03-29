package controllers

import (
	"order-service/common/sentry"
	paymentKafka "order-service/controllers/kafka/payment"
	serviceRegistry "order-service/services"
)

type Registry struct {
	service serviceRegistry.IServiceRegistry
	sentry  sentry.ISentry
}

type IKafkaRegistry interface {
	GetPayment() paymentKafka.IPaymentKafka
}

func NewKafkaRegistry(
	service serviceRegistry.IServiceRegistry,
	sentry sentry.ISentry,
) IKafkaRegistry {
	return &Registry{
		service: service,
		sentry:  sentry,
	}
}

func (r *Registry) GetPayment() paymentKafka.IPaymentKafka {
	return paymentKafka.NewPaymentKafka(r.service, r.sentry)
}
