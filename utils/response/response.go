package response

import (
	"net/http"
	"order-service/utils/sentry"

	constant "order-service/constant/error"
	errorValidation "order-service/utils/error"
)

type Response struct {
	Status  string                               `json:"status"`
	Message any                                  `json:"message"`
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
func ResponseError(err error, sentry sentry.ISentry) Response {
	sentry.CaptureException(err)
	return Response{
		Status:  constant.Error,
		Message: err.Error(),
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
