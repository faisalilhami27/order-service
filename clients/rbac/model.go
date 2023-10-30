package clients

import "github.com/google/uuid"

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
