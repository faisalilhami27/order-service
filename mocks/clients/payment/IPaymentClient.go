// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	context "context"
	clients "order-service/clients/payment"

	mock "github.com/stretchr/testify/mock"
)

// IPaymentClient is an autogenerated mock type for the IPaymentClient type
type IPaymentClient struct {
	mock.Mock
}

// CreatePaymentLink provides a mock function with given fields: _a0, _a1
func (_m *IPaymentClient) CreatePaymentLink(_a0 context.Context, _a1 *clients.PaymentRequest) (*clients.PaymentData, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *clients.PaymentData
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *clients.PaymentRequest) (*clients.PaymentData, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *clients.PaymentRequest) *clients.PaymentData); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*clients.PaymentData)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *clients.PaymentRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIPaymentClient creates a new instance of IPaymentClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIPaymentClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *IPaymentClient {
	mock := &IPaymentClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
