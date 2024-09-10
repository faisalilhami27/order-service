// Code generated by mockery v2.45.1. DO NOT EDIT.

package mocks

import (
	context "context"
	dto "order-service/domain/dto/order"

	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"

	models "order-service/domain/models"

	uuid "github.com/google/uuid"
)

// IOrderRepository is an autogenerated mock type for the IOrderRepository type
type IOrderRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0, _a1, _a2
func (_m *IOrderRepository) Create(_a0 context.Context, _a1 *gorm.DB, _a2 *dto.OrderRequest) (*models.Order, error) {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *models.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *dto.OrderRequest) (*models.Order, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *dto.OrderRequest) *models.Order); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, *dto.OrderRequest) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteByOrderID provides a mock function with given fields: _a0, _a1, _a2
func (_m *IOrderRepository) DeleteByOrderID(_a0 context.Context, _a1 *gorm.DB, _a2 uint) error {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for DeleteByOrderID")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, uint) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindOneOrderByCustomerIDWithLocking provides a mock function with given fields: _a0, _a1, _a2
func (_m *IOrderRepository) FindOneOrderByCustomerIDWithLocking(_a0 context.Context, _a1 *gorm.DB, _a2 uuid.UUID) (*models.Order, error) {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for FindOneOrderByCustomerIDWithLocking")
	}

	var r0 *models.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, uuid.UUID) (*models.Order, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, uuid.UUID) *models.Order); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, uuid.UUID) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOneOrderByID provides a mock function with given fields: _a0, _a1
func (_m *IOrderRepository) FindOneOrderByID(_a0 context.Context, _a1 uint) (*models.Order, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for FindOneOrderByID")
	}

	var r0 *models.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint) (*models.Order, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint) *models.Order); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOneOrderByUUID provides a mock function with given fields: _a0, _a1
func (_m *IOrderRepository) FindOneOrderByUUID(_a0 context.Context, _a1 uuid.UUID) (*models.Order, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for FindOneOrderByUUID")
	}

	var r0 *models.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*models.Order, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *models.Order); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, db, request
func (_m *IOrderRepository) Update(ctx context.Context, db *gorm.DB, request *dto.OrderRequest) error {
	ret := _m.Called(ctx, db, request)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *dto.OrderRequest) error); ok {
		r0 = rf(ctx, db, request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewIOrderRepository creates a new instance of IOrderRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIOrderRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *IOrderRepository {
	mock := &IOrderRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
