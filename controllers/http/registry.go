package controllers

import (
	orderRoute "order-service/controllers/http/suborder"
	serviceRegistry "order-service/services"
	"order-service/utils/sentry"
)

type IControllerRegistry interface {
	GetSubOrder() orderRoute.ISubOrderController
}

type ControllerRegistry struct {
	service serviceRegistry.IServiceRegistry
	sentry  sentry.ISentry
}

func NewControllerRegistry(
	service serviceRegistry.IServiceRegistry,
	sentry sentry.ISentry,
) IControllerRegistry {
	return &ControllerRegistry{
		service: service,
		sentry:  sentry,
	}
}

func (r *ControllerRegistry) GetSubOrder() orderRoute.ISubOrderController {
	return orderRoute.NewOrderController(r.service, r.sentry)
}
