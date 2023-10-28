package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	errorValidation "order-service/utils/error"
	"order-service/utils/response"

	"net/http"

	orderDTO "order-service/domain/dto/suborder"
	"order-service/services"
)

type ISubOrderController interface {
	CreateOrder(c *gin.Context)
	GetSubOrderList(c *gin.Context)
	GetSubOrderDetail(c *gin.Context)
	CancelOrder(c *gin.Context)
}

type ISubOrder struct {
	serviceRegistry services.IServiceRegistry
}

func NewOrderController(serviceRegistry services.IServiceRegistry) ISubOrderController {
	return &ISubOrder{
		serviceRegistry: serviceRegistry,
	}
}

func (o *ISubOrder) CreateOrder(c *gin.Context) {
	var (
		ctx     = c.Request.Context()
		request = orderDTO.SubOrderRequest{}
	)

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError(err))
		return
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		errorResponse := errorValidation.ErrorValidationResponse(err)
		c.JSON(http.StatusUnprocessableEntity, response.ResponseErrorValidation(errorResponse))
		return
	}

	order, err := o.serviceRegistry.GetSubOrder().CreateOrder(ctx, &request)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError(err))
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess(order))
}

func (o *ISubOrder) GetSubOrderList(c *gin.Context) {
	var (
		ctx     = c.Request.Context()
		request = orderDTO.SubOrderRequestParam{}
	)

	err := c.ShouldBindQuery(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError(err))
		return
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		errorResponse := errorValidation.ErrorValidationResponse(err)
		c.JSON(http.StatusUnprocessableEntity, response.ResponseErrorValidation(errorResponse))
		return
	}

	order, err := o.serviceRegistry.GetSubOrder().GetSubOrderList(ctx, &request)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError(err))
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess(order))
}

func (o *ISubOrder) GetSubOrderDetail(c *gin.Context) {
	var (
		ctx       = c.Request.Context()
		orderUUID = c.Param("uuid")
	)

	order, err := o.serviceRegistry.GetSubOrder().GetOrderDetail(ctx, orderUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError(err))
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess(order))
}

func (o *ISubOrder) CancelOrder(c *gin.Context) {
	var (
		ctx       = c.Request.Context()
		orderUUID = c.Param("uuid")
	)

	err := o.serviceRegistry.GetSubOrder().Cancel(ctx, orderUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError(err))
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess(nil))
}
