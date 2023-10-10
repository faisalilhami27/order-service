package utils

import (
	"net/http"
	constant "order-service/constant/error"
)

type Response struct {
	Status  string               `json:"status"`
	Message string               `json:"message"`
	Data    interface{}          `json:"data,omitempty"`
	Error   []ValidationResponse `json:"error,omitempty"`
}

func ResponseSuccess(data interface{}) Response {
	return Response{
		Status:  constant.Success,
		Message: "OK",
		Data:    data,
	}
}

func ResponseError(err error) Response {
	return Response{
		Status:  constant.Error,
		Message: err.Error(),
	}
}

func ResponseErrorValidation(response []ValidationResponse) Response {
	return Response{
		Status:  constant.Error,
		Message: http.StatusText(http.StatusUnprocessableEntity),
		Error:   response,
	}
}
