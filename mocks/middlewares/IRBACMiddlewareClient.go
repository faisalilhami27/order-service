// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	middlewares "order-service/middlewares"

	mock "github.com/stretchr/testify/mock"
)

// IRBACMiddlewareClient is an autogenerated mock type for the IRBACMiddlewareClient type
type IRBACMiddlewareClient struct {
	mock.Mock
}

// CheckPermission provides a mock function with given fields: _a0, _a1
func (_m *IRBACMiddlewareClient) CheckPermission(_a0 string, _a1 []string) (*middlewares.PermissionData, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *middlewares.PermissionData
	var r1 error
	if rf, ok := ret.Get(0).(func(string, []string) (*middlewares.PermissionData, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(string, []string) *middlewares.PermissionData); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*middlewares.PermissionData)
		}
	}

	if rf, ok := ret.Get(1).(func(string, []string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserLogin provides a mock function with given fields: _a0
func (_m *IRBACMiddlewareClient) GetUserLogin(_a0 string) (*middlewares.RBACData, error) {
	ret := _m.Called(_a0)

	var r0 *middlewares.RBACData
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*middlewares.RBACData, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) *middlewares.RBACData); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*middlewares.RBACData)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIRBACMiddlewareClient creates a new instance of IRBACMiddlewareClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIRBACMiddlewareClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *IRBACMiddlewareClient {
	mock := &IRBACMiddlewareClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
