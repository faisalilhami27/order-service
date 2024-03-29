// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	controllers "order-service/controllers/kafka/payment"

	mock "github.com/stretchr/testify/mock"
)

// IKafkaRegistry is an autogenerated mock type for the IKafkaRegistry type
type IKafkaRegistry struct {
	mock.Mock
}

// GetPayment provides a mock function with given fields:
func (_m *IKafkaRegistry) GetPayment() controllers.IPaymentKafka {
	ret := _m.Called()

	var r0 controllers.IPaymentKafka
	if rf, ok := ret.Get(0).(func() controllers.IPaymentKafka); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(controllers.IPaymentKafka)
		}
	}

	return r0
}

// NewIKafkaRegistry creates a new instance of IKafkaRegistry. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIKafkaRegistry(t interface {
	mock.TestingT
	Cleanup(func())
}) *IKafkaRegistry {
	mock := &IKafkaRegistry{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
