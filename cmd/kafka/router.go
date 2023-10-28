package kafka

import (
	"order-service/config"
	kafkaRegistry "order-service/controllers/kafka"
	paymentTopic "order-service/controllers/kafka/payment"

	"golang.org/x/exp/slices"
)

//nolint:revive
type KafkaRouter struct {
	consumer      *ConsumerGroup
	kafkaRegistry kafkaRegistry.IKafkaRegistry
}

type IKafkaRouter interface {
	Register()
}

func NewKafkaRouter(
	consumer *ConsumerGroup,
	registry kafkaRegistry.IKafkaRegistry,
) *KafkaRouter {
	return &KafkaRouter{
		consumer:      consumer,
		kafkaRegistry: registry,
	}
}

func (r *KafkaRouter) Register() {
	r.paymentHandler()
}

func (r *KafkaRouter) paymentHandler() {
	if slices.Contains(config.Config.KafkaConsumerTopics, paymentTopic.PaymentTopic) {
		r.consumer.RegisterTopicHandler(paymentTopic.PaymentTopic, r.kafkaRegistry.GetPayment().HandlePayment)
	}
}
