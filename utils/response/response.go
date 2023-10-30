package response

import (
	"net/http"
	constant "order-service/constant/error"
	errorValidation "order-service/utils/error"
)

type Response struct {
	Status  string                               `json:"status"`
	Message string                               `json:"message"`
	Data    interface{}                          `json:"data,omitempty"`
	Error   []errorValidation.ValidationResponse `json:"error,omitempty"`
}

//nolint:revive
func ResponseSuccess(data interface{}) Response {
	return Response{
		Status:  constant.Success,
		Message: "OK",
		Data:    data,
	}
}

//nolint:revive
func ResponseError(err error) Response {
	return Response{
		Status:  constant.Error,
		Message: err.Error(),
	}
}

//nolint:revive
func ResponsePanicError(err any) Response {
	return Response{
		Status:  constant.Error,
		Message: err.(string), //nolint:forcetypeassert
	}
}

//nolint:revive
func ResponseErrorValidation(response []errorValidation.ValidationResponse) Response {
	return Response{
		Status:  constant.Error,
		Message: http.StatusText(http.StatusUnprocessableEntity),
		Error:   response,
	}
}
