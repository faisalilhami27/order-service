package clients

import (
	"context"
	"fmt"
	"net/http"
	"time"

	clientConfig "order-service/clients/config"
	"order-service/common/sentry"
	"order-service/config"
	"order-service/constant"
	"order-service/utils/helper"
)

type IRbac struct {
	client clientConfig.IClientConfig
	sentry sentry.ISentry
}

type IRbacClient interface {
	GetUserRBAC(context.Context, string) (*RBACData, error)
}

func NewRBACClient(sentry sentry.ISentry, client clientConfig.IClientConfig) IRbacClient {
	return &IRbac{
		client: client,
		sentry: sentry,
	}
}

func (i *IRbac) GetUserRBAC(ctx context.Context, uuid string) (*RBACData, error) {
	logCtx := "common.clients.rbac.rbac.GetUserRBAC"
	var (
		span = i.sentry.StartSpan(ctx, logCtx)
	)
	i.sentry.SpanContext(span)
	defer i.sentry.Finish(span)

	unixTime := time.Now().Unix()
	generateAPIKey := fmt.Sprintf("%s:%s:%d",
		config.Config.AppName,
		i.client.SecretKey(),
		unixTime)
	apiKey := helper.GenerateSHA256(generateAPIKey)

	var response ResponseData
	clone := i.client.Client().Clone().
		Set(constant.XServiceName, config.Config.AppName).
		Set(constant.XApiKey, apiKey).
		Set(constant.XRequestAt, fmt.Sprintf("%d", unixTime)).
		Get(fmt.Sprintf("%s/api/v1/user/%s", i.client.BaseURL(), uuid))

	resp, _, errs := clone.EndStruct(&response)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	if resp.StatusCode != http.StatusOK {
		rbacError := fmt.Errorf("rbac response: %s", response.Message) //nolint:goerr113
		return nil, rbacError
	}

	return &response.Data, nil
}
