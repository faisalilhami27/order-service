package routes

import (
	"github.com/gin-gonic/gin"

	"order-service/controllers"
	"order-service/middlewares"
	subOrderRoute "order-service/routes/suborder"
)

type IRouteRegistry interface {
	Serve()
}

type Route struct {
	controller controllers.IControllerRegistry
	Route      *gin.RouterGroup
}

func NewRouteRegistry(
	controller controllers.IControllerRegistry,
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
