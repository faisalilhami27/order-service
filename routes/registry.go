package routes

import (
	"github.com/gin-gonic/gin"

	"order-service/controllers"
	"order-service/middlewares"
	"order-service/routes/order"
)

type IRouteRegistry interface {
	Serve()
}

type RouteService struct {
	controller controllers.IControllerRegistry
	Route      *gin.RouterGroup
}

func NewRouteRegistry(
	controller controllers.IControllerRegistry,
	route *gin.RouterGroup,
) IRouteRegistry {
	return &RouteService{
		controller: controller,
		Route:      route,
	}
}

func (r *RouteService) Serve() {
	r.Route.Use(middlewares.HandlePanic)
	r.orderRoute().Run()
}

func (r *RouteService) orderRoute() order.IOrderRoute {
	return order.NewOrderRoute(r.controller, r.Route)
}
