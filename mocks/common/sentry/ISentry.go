// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	context "context"

	sentry "github.com/getsentry/sentry-go"
	mock "github.com/stretchr/testify/mock"
)

// ISentry is an autogenerated mock type for the ISentry type
type ISentry struct {
	mock.Mock
}

// CaptureException provides a mock function with given fields: exception
func (_m *ISentry) CaptureException(exception error) *sentry.EventID {
	ret := _m.Called(exception)

	var r0 *sentry.EventID
	if rf, ok := ret.Get(0).(func(error) *sentry.EventID); ok {
		r0 = rf(exception)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sentry.EventID)
		}
	}

	return r0
}

// Finish provides a mock function with given fields: span
func (_m *ISentry) Finish(span *sentry.Span) {
	_m.Called(span)
}

// SpanContext provides a mock function with given fields: span
func (_m *ISentry) SpanContext(span *sentry.Span) context.Context {
	ret := _m.Called(span)

	var r0 context.Context
	if rf, ok := ret.Get(0).(func(*sentry.Span) context.Context); ok {
		r0 = rf(span)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	return r0
}

// StartSpan provides a mock function with given fields: ctx, spanName
func (_m *ISentry) StartSpan(ctx context.Context, spanName string) *sentry.Span {
	ret := _m.Called(ctx, spanName)

	var r0 *sentry.Span
	if rf, ok := ret.Get(0).(func(context.Context, string) *sentry.Span); ok {
		r0 = rf(ctx, spanName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sentry.Span)
		}
	}

	return r0
}

// NewISentry creates a new instance of ISentry. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewISentry(t interface {
	mock.TestingT
	Cleanup(func())
}) *ISentry {
	mock := &ISentry{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}