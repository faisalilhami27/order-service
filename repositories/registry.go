package repositories

import (
	"gorm.io/gorm"

	"order-service/common/sentry"

	orderRepo "order-service/repositories/order"
	orderHistoryRepo "order-service/repositories/orderhistory"
	orderPaymentRepo "order-service/repositories/orderpayment"
	subOrderRepo "order-service/repositories/suborder"
)

type IRepositoryRegistry interface {
	GetTx() *gorm.DB
	GetSubOrder() subOrderRepo.ISubOrderRepository
	GetOrderHistory() orderHistoryRepo.IOrderHistoryRepository
	GetOrderPayment() orderPaymentRepo.IOrderPaymentRepository
	GetOrder() orderRepo.IOrderRepository
}

type Registry struct {
	db     *gorm.DB
	sentry sentry.ISentry
}

func NewRepositoryRegistry(db *gorm.DB, sentry sentry.ISentry) IRepositoryRegistry {
	return &Registry{
		db:     db,
		sentry: sentry,
	}
}

func (r *Registry) GetSubOrder() subOrderRepo.ISubOrderRepository {
	return subOrderRepo.NewSubOrder(r.db, r.sentry)
}

func (r *Registry) GetOrderHistory() orderHistoryRepo.IOrderHistoryRepository {
	return orderHistoryRepo.NewOrderHistory(r.db, r.sentry)
}

func (r *Registry) GetOrderPayment() orderPaymentRepo.IOrderPaymentRepository {
	return orderPaymentRepo.NewOrderPayment(r.db, r.sentry)
}

func (r *Registry) GetOrder() orderRepo.IOrderRepository {
	return orderRepo.NewOrder(r.db, r.sentry)
}

func (r *Registry) GetTx() *gorm.DB {
	return r.db
}
