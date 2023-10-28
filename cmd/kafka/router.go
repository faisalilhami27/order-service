package kafka

import (
	"fmt"
	"order-service/config"
	kafkaRegistry "order-service/controllers/kafka"
	paymentTopic "order-service/controllers/kafka/payment"
	"slices"
)

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
	fmt.Println("payment topic", r.kafkaRegistry.GetPayment().HandlePayment)
	if slices.Contains(config.Config.KafkaConsumerTopics, paymentTopic.PaymentTopic) {
		fmt.Println("payment topic1", paymentTopic.PaymentTopic)
		r.consumer.RegisterTopicHandler(paymentTopic.PaymentTopic, r.kafkaRegistry.GetPayment().HandlePayment)
	}
}
