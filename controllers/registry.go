package controllers

import (
	orderRoute "order-service/controllers/suborder"
	serviceRegistry "order-service/services"
)

type IControllerRegistry interface {
	GetSubOrder() orderRoute.ISubOrderController
}

type ControllerRegistry struct {
	service serviceRegistry.IServiceRegistry
}

func NewControllerRegistry(service serviceRegistry.IServiceRegistry) IControllerRegistry {
	return &ControllerRegistry{
		service: service,
	}
}

func (r *ControllerRegistry) GetSubOrder() orderRoute.ISubOrderController {
	return orderRoute.NewOrderController(r.service)
}
