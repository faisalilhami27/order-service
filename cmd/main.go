package cmd

import (
	"fmt"
	orderPaymentModel "order-service/domain/models/orderpayment"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"order-service/config"
	controllerRegistry "order-service/controllers"
	orderModel "order-service/domain/models/order"
	orderHistoryModel "order-service/domain/models/orderhistory"
	"order-service/migrations"
	repositoryRegistry "order-service/repositories"
	routeRegistry "order-service/routes"
	serviceRegistry "order-service/services"

	"time"
)

var restCmd = &cobra.Command{
	Use:   "serve",
	Short: "Command to start http server",
	Run: func(cmd *cobra.Command, args []string) {
		_ = godotenv.Load() //nolint:errcheck
		config.Init()
		db, err := config.InitDatabase()
		if err != nil {
			panic(err)
		}

		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			panic(err)
		}
		time.Local = loc

		// Database Auto Migration from model
		err = db.AutoMigrate(
			&orderModel.Order{},
			&orderHistoryModel.OrderHistory{},
			&orderPaymentModel.OrderPayment{},
		)
		if err != nil {
			panic(err)
		}

		// Database Auto Migration other than model
		if config.Config.Database.AutoMigrate {
			err = migrations.Run()
			if err != nil {
				panic(err)
			}
		}

		repository := repositoryRegistry.NewRepositoryRegistry(db)
		service := serviceRegistry.NewServiceRegistry(repository)
		controller := controllerRegistry.NewControllerRegistry(service)

		router := gin.Default()
		group := router.Group("/api/v1")
		route := routeRegistry.NewRouteRegistry(controller, group)
		route.Serve()

		port := fmt.Sprintf(":%d", config.Config.Port)
		err = router.Run(port)
		if err != nil {
			panic(err)
		}
	},
}

func Run() {
	err := restCmd.Execute()
	if err != nil {
		panic(err)
	}
}
