package order

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	errorValidation "order-service/utils/error"
	"order-service/utils/response"

	"net/http"

	orderDTO "order-service/domain/dto/order"
	"order-service/services"
)

type IOrderController interface {
	CreateOrder(c *gin.Context)
	GetOrderList(c *gin.Context)
	GetOrderDetail(c *gin.Context)
	CancelOrder(c *gin.Context)
}

type IOrder struct {
	serviceRegistry services.IServiceRegistry
}

func NewOrderController(serviceRegistry services.IServiceRegistry) IOrderController {
	return &IOrder{
		serviceRegistry: serviceRegistry,
	}
}

func (o *IOrder) CreateOrder(c *gin.Context) {
	var (
		ctx     = c.Request.Context()
		request = orderDTO.OrderRequest{}
	)

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError(err))
		return
	}

	order, err := o.serviceRegistry.GetOrder().CreateOrder(ctx, &request)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError(err))
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess(order))
}

func (o *IOrder) GetOrderList(c *gin.Context) {
	var (
		ctx     = c.Request.Context()
		request = orderDTO.OrderRequestParam{}
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

	order, err := o.serviceRegistry.GetOrder().GetOrderList(ctx, &request)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError(err))
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess(order))
}

func (o *IOrder) GetOrderDetail(c *gin.Context) {
	var (
		ctx       = c.Request.Context()
		orderUUID = c.Param("uuid")
	)

	order, err := o.serviceRegistry.GetOrder().GetOrderDetail(ctx, orderUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError(err))
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess(order))
}

func (o *IOrder) CancelOrder(c *gin.Context) {
	var (
		ctx       = c.Request.Context()
		orderUUID = c.Param("uuid")
	)

	err := o.serviceRegistry.GetOrder().Cancel(ctx, orderUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError(err))
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess(nil))
}
