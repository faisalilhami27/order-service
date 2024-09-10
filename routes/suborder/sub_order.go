package routes

import (
	"github.com/gin-gonic/gin"

	controllerRegistry "order-service/controllers/http"
)

type ISubOrderRoute interface {
	Run()
}

type SubOrderRoute struct {
	controller controllerRegistry.IControllerRegistry
	route      *gin.RouterGroup
}

func NewSubOrderRoute(
	controller controllerRegistry.IControllerRegistry,
	route *gin.RouterGroup,
) ISubOrderRoute {
	return &SubOrderRoute{
		controller: controller,
		route:      route,
	}
}

func (o *SubOrderRoute) Run() {
	group := o.route.Group("/order")
	group.GET("", o.controller.GetSubOrder().GetSubOrderList)
	group.GET("/:uuid", o.controller.GetSubOrder().GetSubOrderDetail)
	group.POST("/:uuid", o.controller.GetSubOrder().CancelOrder)
	group.POST("", o.controller.GetSubOrder().CreateOrder)
}
