package routes

import (
	"github.com/gin-gonic/gin"

	"order-service/middlewares"

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
	group.GET("", middlewares.CheckPermission([]string{
		"oms:management-order:order:view",
	}), o.controller.GetSubOrder().GetSubOrderList)
	group.GET("/:uuid", middlewares.CheckPermission([]string{
		"oms:management-order:order:view",
	}), o.controller.GetSubOrder().GetSubOrderDetail)
	group.POST("/:uuid", middlewares.CheckPermission([]string{
		"oms:management-order:order:update",
	}), o.controller.GetSubOrder().CancelOrder)
	group.POST("", middlewares.CheckPermission([]string{
		"oms:management-order:order:create",
	}), o.controller.GetSubOrder().CreateOrder)
}
