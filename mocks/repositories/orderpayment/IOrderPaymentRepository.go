// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	context "context"
	dto "order-service/domain/dto/orderpayment"

	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"

	models "order-service/domain/models/orderpayment"
)

// IOrderPaymentRepository is an autogenerated mock type for the IOrderPaymentRepository type
type IOrderPaymentRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0, _a1, _a2
func (_m *IOrderPaymentRepository) Create(_a0 context.Context, _a1 *gorm.DB, _a2 *dto.OrderPaymentRequest) (*models.OrderPayment, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *models.OrderPayment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *dto.OrderPaymentRequest) (*models.OrderPayment, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *dto.OrderPaymentRequest) *models.OrderPayment); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.OrderPayment)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, *dto.OrderPaymentRequest) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewIOrderPaymentRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewIOrderPaymentRepository creates a new instance of IOrderPaymentRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIOrderPaymentRepository(t mockConstructorTestingTNewIOrderPaymentRepository) *IOrderPaymentRepository {
	mock := &IOrderPaymentRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}