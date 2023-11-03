package cmd

import (
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	clientRegistry "order-service/clients"
	controllerRegistry "order-service/controllers/http"
	kafkaRegistry "order-service/controllers/kafka"
	repositoryRegistry "order-service/repositories"
	routeRegistry "order-service/routes"
	serviceRegistry "order-service/services"

	kafkaConfig "order-service/controllers/kafka/config"

	"order-service/common/circuitbreaker"
	"order-service/common/sentry"
	"order-service/config"
	"order-service/domain/models"
	"order-service/migrations"
	"order-service/utils/response"
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
			&models.Order{},
			&models.SubOrder{},
			&models.OrderHistory{},
			&models.OrderPayment{},
		)
		if err != nil {
			panic(err)
		}

		// Sentry for error tracking
		sentry := sentry.NewSentry(
			sentry.WithDsn(config.Config.SentryDsn),
			sentry.WithDebug(config.Config.AppDebug),
			sentry.WithEnv(config.Config.AppEnv),
			sentry.WithSampleRate(config.Config.SentrySampleRate),
			sentry.WithEnableTracing(config.Config.SentryEnableTracing),
		)

		// Circuit Breaker
		circuitBreaker := circuitbreaker.NewCircuitBreaker(
			sentry,
			circuitbreaker.WithMaxRequest(config.Config.CircuitBreakerMaxRequest),
			circuitbreaker.WithTimeout(config.Config.CircuitBreakerTimeoutInSecond),
		)

		client := clientRegistry.NewClientRegistry(sentry)
		repository := repositoryRegistry.NewRepositoryRegistry(db, sentry)
		service := serviceRegistry.NewServiceRegistry(repository, client, sentry, circuitBreaker)
		controller := controllerRegistry.NewControllerRegistry(service, sentry)

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

		go func() {
			port := fmt.Sprintf(":%d", config.Config.Port)
			err := router.Run(port)
			if err != nil {
				panic(err)
			}
		}()

		// Kafka Consumer
		ctx, cancel := context.WithCancel(context.Background())
		kafkaConsumerConfig := sarama.NewConfig()
		kafkaConsumerConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
			sarama.NewBalanceStrategyRoundRobin()}
		kafkaConsumerConfig.Consumer.Fetch.Default = config.Config.KafkaConsumerFetchDefault
		kafkaConsumerConfig.Consumer.Fetch.Min = config.Config.KafkaConsumerFetchMin
		kafkaConsumerConfig.Consumer.Fetch.Max = config.Config.KafkaConsumerFetchMax
		kafkaConsumerConfig.Consumer.MaxWaitTime = time.Duration(config.Config.KafkaConsumerMaxWaitTimeInMs) * time.Millisecond             //nolint: lll
		kafkaConsumerConfig.Consumer.MaxProcessingTime = time.Duration(config.Config.KafkaConsumerMaxProcessingTimeInMs) * time.Millisecond //nolint: lll
		kafkaConsumerConfig.Consumer.Retry.Backoff = time.Duration(config.Config.KafkaConsumerBackoffTimeInMs) * time.Millisecond           //nolint: lll

		kafkaConsumerClient, err := sarama.NewClient(config.Config.KafkaHosts, kafkaConsumerConfig)
		if err != nil {
			panic(err)
		}
		defer func() {
			if errClose := kafkaConsumerClient.Close(); errClose != nil {
				log.Error(ctx, fmt.Sprintf("error closing update status client: %v", errClose))
			}
		}()

		brokers := config.Config.KafkaHosts
		groupID := config.Config.KafkaConsumerGroupID
		topics := config.Config.KafkaConsumerTopics
		wg := sync.WaitGroup{}
		wg.Add(1)

		if len(topics) > 0 {
			client, err := sarama.NewConsumerGroup(brokers, groupID, kafkaConsumerConfig)
			if err != nil {
				log.Fatal(ctx, fmt.Sprintf("error creating consumer group client: %v", err))
			}

			defer func() {
				if err := client.Close(); err != nil {
					log.Error(ctx, fmt.Sprintf("error closing client: %v", err))
				}
			}()

			consumer := kafkaConfig.NewConsumer()
			kafkaRegistry := kafkaRegistry.NewKafkaRegistry(service, sentry)
			kafkaConsumer := kafkaConfig.NewKafkaRouter(consumer, kafkaRegistry)
			kafkaConsumer.Register()

			KafkaConsumerGroupID, errClient := sarama.NewConsumerGroupFromClient(
				config.Config.KafkaConsumerGroupID,
				kafkaConsumerClient,
			)
			if errClient != nil {
				panic(errClient)
			}

			go func() {
				defer wg.Done()
				defer func() {
					if errClose := KafkaConsumerGroupID.Close(); errClose != nil {
						log.Error(ctx, fmt.Sprintf("error closing update status client: %v", errClose))
					}
				}()

				for {
					errKafkaConsumer := KafkaConsumerGroupID.Consume(ctx, topics, consumer)
					if errKafkaConsumer != nil {
						log.Error(ctx, fmt.Sprintf("error from consumer: %v", errKafkaConsumer))
						return
					}

					if !consumer.KeepRunning() {
						log.Error(ctx, "Consumer is not running anymore")
						return
					}
				}
			}()

			consumer.SetIsReady()

			// Wait for OS signals to gracefully shut down the consumer
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			<-sigChan
		}

		cancel()
		wg.Wait()
	},
}

func Run() {
	err := restCmd.Execute()
	if err != nil {
		panic(err)
	}
}
