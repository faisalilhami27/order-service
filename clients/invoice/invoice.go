package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	clientConfig "order-service/clients/config"
	"order-service/common/sentry"
	"order-service/config"
	"order-service/constant"
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

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	clone := p.client.Client().Clone()
	resp, bodyResp, errs := clone.
		Post(fmt.Sprintf("%s/api/v1/invoice/generate", p.client.BaseURL())).
		Set(constant.XApiKey, config.Config.InternalService.Invoice.StaticKey).
		Send(string(body)).
		End()

	var errResponse ErrorInvoiceResponse
	if resp.StatusCode != http.StatusOK || len(errs) > 0 {
		err = json.Unmarshal([]byte(bodyResp), &errResponse)
		if err != nil {
			return nil, err
		}
		invoiceError := fmt.Errorf("invoice response: %s", errResponse.Message) //nolint:goerr113
		return nil, invoiceError
	}

	var response InvoiceResponse
	err = json.Unmarshal([]byte(bodyResp), &response)
	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}
