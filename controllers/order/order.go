package order

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"net/http"

	orderDTO "order-service/domain/dto/order"
	"order-service/services"
	"order-service/utils"
)

type IOrderController interface {
	CreateOrder(c *gin.Context)
	GetOrderList(c *gin.Context)
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
		c.JSON(http.StatusBadRequest, utils.ResponseError(err))
		return
	}

	order, err := o.serviceRegistry.GetOrder().CreateOrder(ctx, &request)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseSuccess(order))
}

func (o *IOrder) GetOrderList(c *gin.Context) {
	var (
		ctx     = c.Request.Context()
		request = orderDTO.OrderRequestParam{}
	)

	err := c.ShouldBindQuery(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err))
		return
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		errorResponse := utils.ErrorResponse(err)
		c.JSON(http.StatusUnprocessableEntity, utils.ResponseErrorValidation(errorResponse))
		return
	}

	order, err := o.serviceRegistry.GetOrder().GetOrderList(ctx, &request)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseSuccess(order))
}
