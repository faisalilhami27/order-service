package clients

import (
	"encoding/json"
	"fmt"

	"net/http"

	clientConfig "order-service/clients/config"
	"order-service/config"
	"order-service/constant"
	"order-service/utils/helper"

	"time"
)

type IPayment struct {
	client clientConfig.IClientConfig
}

type IPaymentClient interface {
	CreatePaymentLink(*PaymentRequest) (*PaymentData, error)
}

func NewPaymentClient(client clientConfig.IClientConfig) IPaymentClient {
	return &IPayment{
		client: client,
	}
}

func (p *IPayment) CreatePaymentLink(request *PaymentRequest) (*PaymentData, error) {
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

	var errResponse ErrorResponse
	if resp.StatusCode != http.StatusCreated || len(errs) > 0 {
		err = json.Unmarshal([]byte(bodyResp), &errResponse)
		if err != nil {
			return nil, err
		}
		paymentError := fmt.Errorf("payment response: %s", errResponse.Message) //nolint:goerr113
		return nil, paymentError
	}

	var response PaymentResponse
	err = json.Unmarshal([]byte(bodyResp), &response)
	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}
