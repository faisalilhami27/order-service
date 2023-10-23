package routes

import (
	"github.com/gin-gonic/gin"

	"order-service/controllers"
)

type ISubOrderRoute interface {
	Run()
}

type SubOrderRoute struct {
	controller controllers.IControllerRegistry
	route      *gin.RouterGroup
}

func NewSubOrderRoute(
	controller controllers.IControllerRegistry,
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
