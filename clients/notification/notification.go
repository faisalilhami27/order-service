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

type INotification struct {
	client clientConfig.IClientConfig
	sentry sentry.ISentry
}

type INotificationClient interface {
	SendToWhatsapp(context.Context, *NotificationRequest) error
}

func NewNotificationClient(
	sentry sentry.ISentry,
	client clientConfig.IClientConfig,
) INotificationClient {
	return &INotification{
		client: client,
		sentry: sentry,
	}
}

func (p *INotification) SendToWhatsapp(ctx context.Context, request *NotificationRequest) error {
	logCtx := "common.clients.notification.notification.SendToWhatsapp"
	var (
		span = p.sentry.StartSpan(ctx, logCtx)
	)
	p.sentry.SpanContext(span)
	defer p.sentry.Finish(span)

	body, err := json.Marshal(request)
	if err != nil {
		return err
	}

	clone := p.client.Client().Clone()
	resp, bodyResp, errs := clone.
		Post(fmt.Sprintf("%s/api/v1/template/send-message", p.client.BaseURL())).
		Set(constant.XApiKey, config.Config.InternalService.Notification.StaticKey).
		Send(string(body)).
		End()

	if len(errs) > 0 {
		return errs[0]
	}

	var errResponse ErrorNotificationResponse
	if resp.StatusCode != http.StatusOK {
		err = json.Unmarshal([]byte(bodyResp), &errResponse)
		if err != nil {
			return err
		}
		notificationError := fmt.Errorf("notification response: %s", errResponse.Message) //nolint:goerr113
		return notificationError
	}

	var response NotificationResponse
	err = json.Unmarshal([]byte(bodyResp), &response)
	if err != nil {
		return err
	}

	return nil
}
