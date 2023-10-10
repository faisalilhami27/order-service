package order

import (
	"github.com/gin-gonic/gin"

	"order-service/controllers"
)

type IOrderRoute interface {
	Run()
}

type OrderRoute struct { //nolint:revive
	controller controllers.IControllerRegistry
	route      *gin.RouterGroup
}

func NewOrderRoute(
	controller controllers.IControllerRegistry,
	route *gin.RouterGroup,
) IOrderRoute {
	return &OrderRoute{
		controller: controller,
		route:      route,
	}
}

func (o *OrderRoute) Run() {
	group := o.route.Group("/order")
	group.GET("", o.controller.GetOrder().GetOrderList)
	group.POST("", o.controller.GetOrder().CreateOrder)
}
