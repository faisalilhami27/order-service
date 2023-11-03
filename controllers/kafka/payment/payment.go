package controllers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"order-service/common/sentry"
	errorResp "order-service/utils/error"

	dto "order-service/domain/dto/kafka/payment"
	paymentDTO "order-service/domain/dto/suborder"
	serviceRegistry "order-service/services"
)

const PaymentTopic = "payment-service-callback"

type PaymentKafka struct {
	service serviceRegistry.IServiceRegistry
	sentry  sentry.ISentry
}

type IPaymentKafka interface {
	HandlePayment(ctx context.Context, message *sarama.ConsumerMessage) error
}

func NewPaymentKafka(
	service serviceRegistry.IServiceRegistry,
	sentry sentry.ISentry,
) IPaymentKafka {
	return &PaymentKafka{
		service: service,
		sentry:  sentry,
	}
}

func (p *PaymentKafka) HandlePayment(ctx context.Context, message *sarama.ConsumerMessage) error {
	var body dto.PaymentContent
	err := json.Unmarshal(message.Value, &body)
	if err != nil {
		log.Errorf(fmt.Sprintf("error unmarshal: %s", err.Error()), err)
		return err
	}

	data := body.Body.Data
	orderUUID, _ := uuid.Parse(data.OrderID)     //nolint:errcheck
	paymentUUID, _ := uuid.Parse(data.PaymentID) //nolint:errcheck
	switch body.Event.Name {
	case "PENDING":
		err = p.service.GetSubOrder().ReceivePendingPayment(ctx, &paymentDTO.PaymentRequest{
			OrderID:     orderUUID,
			PaymentID:   paymentUUID,
			PaymentLink: data.PaymentLink,
			PaymentType: data.PaymentType,
			Amount:      data.Amount,
			Status:      data.Status,
			VaNumber:    data.VANumber,
			Bank:        data.Bank,
			Acquirer:    data.Acquirer,
		})
	case "SETTLEMENT":
		err = p.service.GetSubOrder().ReceivePaymentSettlement(ctx, &paymentDTO.PaymentRequest{
			OrderID:     orderUUID,
			PaymentID:   paymentUUID,
			PaymentLink: data.PaymentLink,
			PaymentType: data.PaymentType,
			Amount:      data.Amount,
			Status:      data.Status,
			VaNumber:    data.VANumber,
			Bank:        data.Bank,
			Acquirer:    data.Acquirer,
			PaidAt:      data.PaidAt,
		})
	case "EXPIRE":
		err = p.service.GetSubOrder().ReceivePaymentExpire(ctx, &paymentDTO.PaymentRequest{
			OrderID:   orderUUID,
			PaymentID: paymentUUID,
			Status:    data.Status,
		})
	}
	if err != nil {
		return errorResp.WrapError(err, p.sentry)
	}
	return nil
}
