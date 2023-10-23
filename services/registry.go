package services

import (
	"order-service/clients"
	repositoryRegistry "order-service/repositories"
	orderService "order-service/services/suborder"
)

type IServiceRegistry interface {
	GetSubOrder() orderService.ISubOrderService
}

type Registry struct {
	repository repositoryRegistry.IRepositoryRegistry
	client     clients.IClientRegistry
}

func NewServiceRegistry(
	repository repositoryRegistry.IRepositoryRegistry,
	client clients.IClientRegistry,
) IServiceRegistry {
	return &Registry{
		repository: repository,
		client:     client,
	}
}

func (s *Registry) GetSubOrder() orderService.ISubOrderService {
	return orderService.NewSubOrderService(s.repository, s.client)
}
