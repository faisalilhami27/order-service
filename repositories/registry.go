package repositories

import (
	"gorm.io/gorm"

	orderRepo "order-service/repositories/order"
	orderHistoryRepo "order-service/repositories/orderhistory"
	orderPaymentRepo "order-service/repositories/orderpayment"
)

type IRepositoryRegistry interface {
	GetTx() *gorm.DB
	GetOrderRepository() orderRepo.IOrderRepository
	GetOrderHistoryRepository() orderHistoryRepo.IOrderHistoryRepository
	GetOrderPaymentRepository() orderPaymentRepo.IOrderPaymentRepository
}

type Registry struct {
	db *gorm.DB
}

func NewRepositoryRegistry(db *gorm.DB) IRepositoryRegistry {
	return &Registry{
		db: db,
	}
}

func (r *Registry) GetOrderRepository() orderRepo.IOrderRepository {
	return orderRepo.NewOrder(r.db)
}

func (r *Registry) GetOrderHistoryRepository() orderHistoryRepo.IOrderHistoryRepository {
	return orderHistoryRepo.NewOrderHistory(r.db)
}

func (r *Registry) GetOrderPaymentRepository() orderPaymentRepo.IOrderPaymentRepository {
	return orderPaymentRepo.NewOrderPayment(r.db)
}

func (r *Registry) GetTx() *gorm.DB {
	return r.db
}
