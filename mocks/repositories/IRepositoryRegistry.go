// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	orderhistory "order-service/repositories/orderhistory"

	mock "github.com/stretchr/testify/mock"
	gorm "gorm.io/gorm"

	orderpayment "order-service/repositories/orderpayment"

	suborder "order-service/repositories/suborder"
)

// IRepositoryRegistry is an autogenerated mock type for the IRepositoryRegistry type
type IRepositoryRegistry struct {
	mock.Mock
}

// GetOrderHistoryRepository provides a mock function with given fields:
func (_m *IRepositoryRegistry) GetOrderHistoryRepository() orderhistory.IOrderHistoryRepository {
	ret := _m.Called()

	var r0 orderhistory.IOrderHistoryRepository
	if rf, ok := ret.Get(0).(func() orderhistory.IOrderHistoryRepository); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(orderhistory.IOrderHistoryRepository)
		}
	}

	return r0
}

// GetOrderPaymentRepository provides a mock function with given fields:
func (_m *IRepositoryRegistry) GetOrderPaymentRepository() orderpayment.IOrderPaymentRepository {
	ret := _m.Called()

	var r0 orderpayment.IOrderPaymentRepository
	if rf, ok := ret.Get(0).(func() orderpayment.IOrderPaymentRepository); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(orderpayment.IOrderPaymentRepository)
		}
	}

	return r0
}

// GetSubOrderRepository provides a mock function with given fields:
func (_m *IRepositoryRegistry) GetSubOrderRepository() suborder.ISubOrderRepository {
	ret := _m.Called()

	var r0 suborder.ISubOrderRepository
	if rf, ok := ret.Get(0).(func() suborder.ISubOrderRepository); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(suborder.ISubOrderRepository)
		}
	}

	return r0
}

// GetTx provides a mock function with given fields:
func (_m *IRepositoryRegistry) GetTx() *gorm.DB {
	ret := _m.Called()

	var r0 *gorm.DB
	if rf, ok := ret.Get(0).(func() *gorm.DB); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorm.DB)
		}
	}

	return r0
}

type mockConstructorTestingTNewIRepositoryRegistry interface {
	mock.TestingT
	Cleanup(func())
}

// NewIRepositoryRegistry creates a new instance of IRepositoryRegistry. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIRepositoryRegistry(t mockConstructorTestingTNewIRepositoryRegistry) *IRepositoryRegistry {
	mock := &IRepositoryRegistry{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
