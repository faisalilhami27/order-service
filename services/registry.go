package services

import (
	"order-service/clients"
	"order-service/common/circuitbreaker"
	"order-service/common/sentry"
	repositoryRegistry "order-service/repositories"
	orderService "order-service/services/suborder"
)

type IServiceRegistry interface {
	GetSubOrder() orderService.ISubOrderService
}

type Registry struct {
	repository repositoryRegistry.IRepositoryRegistry
	client     clients.IClientRegistry
	sentry     sentry.ISentry
	breaker    circuitbreaker.ICircuitBreaker
}

func NewServiceRegistry(
	repository repositoryRegistry.IRepositoryRegistry,
	client clients.IClientRegistry,
	sentry sentry.ISentry,
	breaker circuitbreaker.ICircuitBreaker,
) IServiceRegistry {
	return &Registry{
		repository: repository,
		client:     client,
		sentry:     sentry,
		breaker:    breaker,
	}
}

func (s *Registry) GetSubOrder() orderService.ISubOrderService {
	return orderService.NewSubOrderService(s.repository, s.client, s.sentry, s.breaker)
}
