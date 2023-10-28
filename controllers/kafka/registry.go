package controllers

import (
	paymentKafka "order-service/controllers/kafka/payment"
	serviceRegistry "order-service/services"
)

type Registry struct {
	service serviceRegistry.IServiceRegistry
}

type IKafkaRegistry interface {
	GetPayment() paymentKafka.IPaymentKafka
}

func NewKafkaRegistry(service serviceRegistry.IServiceRegistry) IKafkaRegistry {
	return &Registry{
		service: service,
	}
}

func (r *Registry) GetPayment() paymentKafka.IPaymentKafka {
	return paymentKafka.NewPaymentKafka(r.service)
}
