package services

import (
	repositoryRegistry "order-service/repositories"
	orderService "order-service/services/order"
)

type IServiceRegistry interface {
	GetOrder() orderService.IOrderService
}

type Registry struct {
	repository repositoryRegistry.IRepositoryRegistry
}

func NewServiceRegistry(repository repositoryRegistry.IRepositoryRegistry) IServiceRegistry {
	return &Registry{
		repository: repository,
	}
}

func (s *Registry) GetOrder() orderService.IOrderService {
	return orderService.NewOrderService(s.repository)
}
