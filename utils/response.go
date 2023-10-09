package utils

import (
	constant "order-service/constant/error"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func ResponseSuccess(data interface{}) Response {
	return Response{
		Status:  constant.Success,
		Message: constant.Success,
		Data:    data,
	}
}

func ResponseError(err error) Response {
	return Response{
		Status:  constant.Error,
		Message: err.Error(),
	}
}
