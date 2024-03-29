// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	gin "github.com/gin-gonic/gin"
	mock "github.com/stretchr/testify/mock"
)

// ISubOrderController is an autogenerated mock type for the ISubOrderController type
type ISubOrderController struct {
	mock.Mock
}

// CancelOrder provides a mock function with given fields: c
func (_m *ISubOrderController) CancelOrder(c *gin.Context) {
	_m.Called(c)
}

// CreateOrder provides a mock function with given fields: c
func (_m *ISubOrderController) CreateOrder(c *gin.Context) {
	_m.Called(c)
}

// GetSubOrderDetail provides a mock function with given fields: c
func (_m *ISubOrderController) GetSubOrderDetail(c *gin.Context) {
	_m.Called(c)
}

// GetSubOrderList provides a mock function with given fields: c
func (_m *ISubOrderController) GetSubOrderList(c *gin.Context) {
	_m.Called(c)
}

// NewISubOrderController creates a new instance of ISubOrderController. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewISubOrderController(t interface {
	mock.TestingT
	Cleanup(func())
}) *ISubOrderController {
	mock := &ISubOrderController{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
