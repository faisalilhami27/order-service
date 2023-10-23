package cmd

import (
	"fmt"
	"net/http"
	"order-service/utils/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	orderModel "order-service/domain/models/order"
	orderHistoryModel "order-service/domain/models/orderhistory"
	orderPaymentModel "order-service/domain/models/orderpayment"
	subOrderModel "order-service/domain/models/suborder"

	"order-service/config"
	"order-service/migrations"

	clientRegistry "order-service/clients"
	controllerRegistry "order-service/controllers"
	repositoryRegistry "order-service/repositories"
	routeRegistry "order-service/routes"
	serviceRegistry "order-service/services"
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

		// Database Auto Migration other than model
		if config.Config.Database.AutoMigrate {
			err = migrations.Run()
			if err != nil {
				panic(err)
			}
		}

		// Database Auto Migration from model
		err = db.AutoMigrate(
			&orderModel.Order{},
			&subOrderModel.SubOrder{},
			&orderHistoryModel.OrderHistory{},
			&orderPaymentModel.OrderPayment{},
		)
		if err != nil {
			panic(err)
		}

		client := clientRegistry.NewClientRegistry()
		repository := repositoryRegistry.NewRepositoryRegistry(db)
		service := serviceRegistry.NewServiceRegistry(repository, client)
		controller := controllerRegistry.NewControllerRegistry(service)

		router := gin.Default()
		router.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, response.Response{
				Status:  "error",
				Message: fmt.Sprintf("Path %s", http.StatusText(http.StatusNotFound)),
			})
		})

		router.Use(func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Next()
		})
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
