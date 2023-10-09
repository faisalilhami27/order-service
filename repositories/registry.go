package repositories

import (
	"gorm.io/gorm"

	orderRepo "order-service/repositories/order"
	"order-service/repositories/orderhistory"
)

type IRepositoryRegistry interface {
	GetTx() *gorm.DB
	GetOrderRepository() orderRepo.IOrderRepository
	GetOrderHistoryRepository() orderhistory.IOrderHistoryRepository
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

func (r *Registry) GetOrderHistoryRepository() orderhistory.IOrderHistoryRepository {
	return orderhistory.NewOrder(r.db)
}

func (r *Registry) GetTx() *gorm.DB {
	return r.db
}
