package services

import (
	"order-service/clients"
	repositoryRegistry "order-service/repositories"
	orderService "order-service/services/suborder"
	"order-service/utils/sentry"
)

type IServiceRegistry interface {
	GetSubOrder() orderService.ISubOrderService
}

type Registry struct {
	repository repositoryRegistry.IRepositoryRegistry
	client     clients.IClientRegistry
	sentry     sentry.ISentry
}

func NewServiceRegistry(
	repository repositoryRegistry.IRepositoryRegistry,
	client clients.IClientRegistry,
	sentry sentry.ISentry,
) IServiceRegistry {
	return &Registry{
		repository: repository,
		client:     client,
		sentry:     sentry,
	}
}

func (s *Registry) GetSubOrder() orderService.ISubOrderService {
	return orderService.NewSubOrderService(s.repository, s.client, s.sentry)
}
