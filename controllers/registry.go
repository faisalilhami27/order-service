package controllers

import (
	orderRoute "order-service/controllers/order"
	serviceRegistry "order-service/services"
)

type IControllerRegistry interface {
	GetOrder() orderRoute.IOrderController
}

type ControllerRegistry struct {
	service serviceRegistry.IServiceRegistry
}

func NewControllerRegistry(service serviceRegistry.IServiceRegistry) IControllerRegistry {
	return &ControllerRegistry{
		service: service,
	}
}

func (r *ControllerRegistry) GetOrder() orderRoute.IOrderController {
	return orderRoute.NewOrderController(r.service)
}
