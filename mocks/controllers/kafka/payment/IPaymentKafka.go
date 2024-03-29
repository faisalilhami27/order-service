// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	sarama "github.com/IBM/sarama"
)

// IPaymentKafka is an autogenerated mock type for the IPaymentKafka type
type IPaymentKafka struct {
	mock.Mock
}

// HandlePayment provides a mock function with given fields: ctx, message
func (_m *IPaymentKafka) HandlePayment(ctx context.Context, message *sarama.ConsumerMessage) error {
	ret := _m.Called(ctx, message)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *sarama.ConsumerMessage) error); ok {
		r0 = rf(ctx, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewIPaymentKafka creates a new instance of IPaymentKafka. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIPaymentKafka(t interface {
	mock.TestingT
	Cleanup(func())
}) *IPaymentKafka {
	mock := &IPaymentKafka{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
