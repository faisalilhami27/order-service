package order

import (
	"github.com/gin-gonic/gin"

	"net/http"

	orderDTO "order-service/domain/dto/order"
	"order-service/services"
	"order-service/utils"
)

type IOrderController interface {
	CreateOrder(c *gin.Context)
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
