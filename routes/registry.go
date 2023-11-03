package routes

import (
	"github.com/gin-gonic/gin"

	controllerRegistry "order-service/controllers/http"
	"order-service/middlewares"
	subOrderRoute "order-service/routes/suborder"
)

type IRouteRegistry interface {
	Serve()
}

type Route struct {
	controller controllerRegistry.IControllerRegistry
	Route      *gin.RouterGroup
}

func NewRouteRegistry(
	controller controllerRegistry.IControllerRegistry,
	route *gin.RouterGroup,
) IRouteRegistry {
	return &Route{
		controller: controller,
		Route:      route,
	}
}

func (r *Route) Serve() {
	r.Route.Use(middlewares.HandlePanic)
	r.suOrderRoute().Run()
}

func (r *Route) suOrderRoute() subOrderRoute.ISubOrderRoute {
	return subOrderRoute.NewSubOrderRoute(r.controller, r.Route)
}
