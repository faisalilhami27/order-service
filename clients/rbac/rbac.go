package clients

import (
	"fmt"

	"net/http"

	clientConfig "order-service/clients/config"
	"order-service/config"
	"order-service/constant"
	"order-service/utils/helper"

	"time"
)

type IRbac struct {
	client clientConfig.IClientConfig
}

type IRbacClient interface {
	GetUserRBAC(string) (*RBACData, error)
}

func NewRBACClient(client clientConfig.IClientConfig) IRbacClient {
	return &IRbac{
		client: client,
	}
}

func (i *IRbac) GetUserRBAC(uuid string) (*RBACData, error) {
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
