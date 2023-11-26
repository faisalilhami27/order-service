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

type IInvoice struct {
	client clientConfig.IClientConfig
	sentry sentry.ISentry
}

type IInvoiceClient interface {
	GenerateInvoice(context.Context, *InvoiceRequest) (*InvoiceData, error)
}

func NewInvoiceClient(
	sentry sentry.ISentry,
	client clientConfig.IClientConfig,
) IInvoiceClient {
	return &IInvoice{
		client: client,
		sentry: sentry,
	}
}

func (p *IInvoice) GenerateInvoice(ctx context.Context, request *InvoiceRequest) (*InvoiceData, error) {
	logCtx := "common.clients.invoice.invoice.GenerateInvoice"
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
		Post(fmt.Sprintf("%s/api/v1/invoice/generate", p.client.BaseURL())).
		Set(constant.XServiceName, config.Config.AppName).
		Set(constant.XApiKey, apiKey).
		Set(constant.XRequestAt, fmt.Sprintf("%d", unixTime)).
		Send(string(body)).
		End()

	var errResponse ErrorInvoiceResponse
	if resp.StatusCode != http.StatusCreated || len(errs) > 0 {
		err = json.Unmarshal([]byte(bodyResp), &errResponse)
		if err != nil {
			return nil, err
		}
		paymentError := fmt.Errorf("invoice response: %s", errResponse.Message) //nolint:goerr113
		return nil, paymentError
	}

	var response InvoiceResponse
	err = json.Unmarshal([]byte(bodyResp), &response)
	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}
