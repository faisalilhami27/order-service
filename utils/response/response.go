package response

import (
	"github.com/gin-gonic/gin"

	"net/http"

	constant "order-service/constant/error"
	"order-service/utils/sentry"
)

type Response struct {
	Status  string      `json:"status"`
	Message any         `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ParamHTTPResp struct {
	Code    int
	Err     error
	Message *string
	Gin     *gin.Context
	Data    interface{}
	Sentry  sentry.ISentry
}

func HTTPResponse(param ParamHTTPResp) {
	if param.Err == nil {
		param.Gin.JSON(param.Code, Response{
			Status:  constant.Success,
			Message: http.StatusText(http.StatusOK),
			Data:    param.Data,
		})
		return
	}

	var message = param.Err.Error()
	if param.Message != nil {
		message = *param.Message
	}

	param.Gin.JSON(param.Code, Response{
		Status:  constant.Error,
		Message: message,
		Data:    param.Data,
	})
	param.Sentry.CaptureException(param.Err)
	return //nolint:gosimple
}
