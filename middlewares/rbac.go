package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	clientConfig "order-service/clients/config"
	"order-service/config"
	"order-service/constant"
	"order-service/utils/helper"
)

type ErrorRBACResponse struct {
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Data    *interface{} `json:"data,omitempty"`
}

type RBACResponse[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

type ResponseData struct {
	RBACResponse[RBACData]
}

type RBACData struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	PhoneNumber string    `json:"phone_number"`
	Roles       []Entity  `json:"roles"`
	Permissions []Entity  `json:"permissions"`
}

type Entity struct {
	Name string `json:"name"`
}

type PermissionData struct {
	Allowed        bool     `json:"allowed"`
	MissPermission []string `json:"miss_permission"`
}

type IRBACMiddleware struct {
	client clientConfig.IClientConfig
}

type PermissionRequest struct {
	Permissions []string `json:"permissions"`
}

type IRBACMiddlewareClient interface {
	GetUserLogin(string) (*RBACData, error)
	CheckPermission(string, []string) (*PermissionData, error)
}

func NewRBACMiddleware(client clientConfig.IClientConfig) IRBACMiddlewareClient {
	return &IRBACMiddleware{
		client: client,
	}
}

func (m *IRBACMiddleware) GetUserLogin(token string) (*RBACData, error) {
	unixTime := time.Now().Unix()
	generateAPIKey := fmt.Sprintf("%s:%s:%d",
		config.Config.AppName,
		m.client.SecretKey(),
		unixTime)
	apiKey := helper.GenerateSHA256(generateAPIKey)

	var response ResponseData
	clone := m.client.Client().Clone().
		Set(constant.XServiceName, config.Config.AppName).
		Set(constant.XApiKey, apiKey).
		Set(constant.XRequestAt, fmt.Sprintf("%d", unixTime)).
		Set(constant.Authorization, token).
		Get(fmt.Sprintf("%s/api/v1/user/login", m.client.BaseURL()))

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

func (m *IRBACMiddleware) CheckPermission(token string, permissions []string) (*PermissionData, error) {
	unixTime := time.Now().Unix()
	generateAPIKey := fmt.Sprintf("%s:%s:%d",
		config.Config.AppName,
		m.client.SecretKey(),
		unixTime)
	apiKey := helper.GenerateSHA256(generateAPIKey)

	request := PermissionRequest{
		Permissions: permissions,
	}
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	clone := m.client.Client().Clone()
	resp, bodyResp, errs := clone.
		Post(fmt.Sprintf("%s/api/v1/permission/check", m.client.BaseURL())).
		Set(constant.XServiceName, config.Config.AppName).
		Set(constant.XApiKey, apiKey).
		Set(constant.XRequestAt, fmt.Sprintf("%d", unixTime)).
		Set(constant.Authorization, token).
		Send(string(body)).
		End()

	var errResponse ErrorRBACResponse
	if resp.StatusCode != http.StatusOK || len(errs) > 0 {
		err = json.Unmarshal([]byte(bodyResp), &errResponse)
		if err != nil {
			return nil, err
		}
		rbacError := fmt.Errorf("rbac response: %s", errResponse.Message) //nolint:goerr113
		return nil, rbacError
	}

	var response RBACResponse[PermissionData]
	err = json.Unmarshal([]byte(bodyResp), &response)
	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}
