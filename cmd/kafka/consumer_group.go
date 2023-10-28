package kafka

import (
	"context"
	"fmt"
	"order-service/config"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"order-service/constant"
)

type (
	Payload struct {
		Key     []byte
		Headers map[string][]byte
		Value   []byte
	}
	TopicName string
	Handler   func(ctx context.Context, message *sarama.ConsumerMessage) error
)

type ConsumerGroup struct {
	mu          *sync.Mutex
	isReady     chan bool
	keepRunning bool
	handlers    map[TopicName]Handler
	retryWg     sync.WaitGroup
}

func NewConsumer() *ConsumerGroup {
	return &ConsumerGroup{
		mu:          &sync.Mutex{},
		isReady:     make(chan bool),
		keepRunning: true,
		handlers:    make(map[TopicName]Handler),
	}
}

func (c *ConsumerGroup) Setup(sarama.ConsumerGroupSession) error {
	log.Infof("ConsumerGroup setup done")
	close(c.isReady)
	return nil
}

func (c *ConsumerGroup) Cleanup(sarama.ConsumerGroupSession) error {
	log.Infof("ConsumerGroup cleanup")
	return nil
}

//nolint:gocognit,cyclop
func (c *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	messageChan := claim.Messages()

	for {
		select {
		case message, ok := <-messageChan:
			if !ok {
				return nil
			}

			handler, exists := c.handlers[TopicName(message.Topic)]
			if !exists {
				log.Errorf("No handler found for topic: %s", message.Topic)
				continue
			}

			ctx := context.Background()

			var requestID string
			if ctx.Value(constant.XRequestID) != nil {
				requestID, isValid := ctx.Value(constant.XRequestID).(string)
				if !isValid {
					log.Errorf("invalid request id")
				}

				_, err := uuid.Parse(requestID)
				if err != nil {
					log.Errorf("uuid parse error: %v", err)
					ctx = c.generateRequestID(ctx, nil)
				}
			} else {
				requestID = uuid.New().String()
				ctx = c.generateRequestID(ctx, &requestID)
			}

			var err error
			retries := config.Config.KafkaMaxRetry

			for retries > 0 {
				c.retryWg.Add(1)
				go func(retries int) {
					defer c.retryWg.Done()
					err = handler(ctx, message)
				}(retries)

				c.retryWg.Wait()

				if err == nil {
					break
				}

				log.Errorf("Error handling message: %v. Retrying...", err)
				retries--
			}

			if err != nil {
				log.Errorf("Error handling message after retries: %v", err)
			}

			session.MarkMessage(message, time.Now().UTC().String())

		case <-session.Context().Done():
			return nil
		}
	}
}

func (c *ConsumerGroup) SetIsReady() {
	<-c.isReady
}

func (c *ConsumerGroup) KeepRunning() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.keepRunning
}

func (c *ConsumerGroup) Recover() {
	if r := recover(); r != nil {
		log.Errorf("Recovered from panic: %v", r)
	}
}

func (c *ConsumerGroup) RegisterTopicHandler(topicName TopicName, handler Handler) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.handlers[topicName] = handler
	log.Infof(fmt.Sprintf("registered handler: %s", topicName))
}

func (c *ConsumerGroup) generateRequestID(ctx context.Context, requestID *string) context.Context {
	if requestID == nil {
		uuid := uuid.New().String()
		requestID = &uuid
	}
	ctx = context.WithValue(
		ctx,
		constant.XRequestID, //nolint:staticcheck
		*requestID,
	)

	return ctx
}
