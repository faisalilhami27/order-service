package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	clientConfig "order-service/clients/config"
	"order-service/common/sentry"
	"order-service/config"
	"order-service/constant"
	"order-service/utils/helper"
)

type IPayment struct {
	client clientConfig.IClientConfig
	sentry sentry.ISentry
}

type IPaymentClient interface {
	CreatePaymentLink(context.Context, *PaymentRequest) (*PaymentData, error)
}

func NewPaymentClient(
	sentry sentry.ISentry,
	client clientConfig.IClientConfig,
) IPaymentClient {
	return &IPayment{
		client: client,
		sentry: sentry,
	}
}

func (p *IPayment) CreatePaymentLink(ctx context.Context, request *PaymentRequest) (*PaymentData, error) {
	logCtx := "common.clients.payment.payment.CreatePaymentLink"
	var (
		span = p.sentry.StartSpan(ctx, logCtx)
	)
	p.sentry.SpanContext(span)
	defer p.sentry.Finish(span)

	unixTime := time.Now().Unix()
	generateAPIKey := fmt.Sprintf("%s:%s:%d",
		config.Config.AppName,
		p.client.SecretKey(),
		unixTime)
	apiKey := helper.GenerateSHA256(generateAPIKey)

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	clone := p.client.Client().Clone()
	resp, bodyResp, errs := clone.
		Post(fmt.Sprintf("%s/api/v1/payment", p.client.BaseURL())).
		Set(constant.XServiceName, config.Config.AppName).
		Set(constant.XApiKey, apiKey).
		Set(constant.XRequestAt, fmt.Sprintf("%d", unixTime)).
		Send(string(body)).
		End()

	if len(errs) > 0 {
		return nil, errs[0]
	}

	var errResponse ErrorPaymentResponse
	if resp.StatusCode != http.StatusCreated {
		err = json.Unmarshal([]byte(bodyResp), &errResponse)
		if err != nil {
			return nil, err
		}
		paymentError := fmt.Errorf("payment response: %s", errResponse.Message) //nolint:err113
		return nil, paymentError
	}

	var response PaymentResponse
	err = json.Unmarshal([]byte(bodyResp), &response)
	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}
