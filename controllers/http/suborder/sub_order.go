package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	errorValidation "order-service/utils/error"
	"order-service/utils/response"
	"order-service/utils/sentry"

	"net/http"

	orderDTO "order-service/domain/dto/suborder"
	"order-service/services"
)

type ISubOrderController interface {
	CreateOrder(c *gin.Context)
	GetSubOrderList(c *gin.Context)
	GetSubOrderDetail(c *gin.Context)
	CancelOrder(c *gin.Context)
}

type ISubOrder struct {
	serviceRegistry services.IServiceRegistry
	sentry          sentry.ISentry
}

func NewOrderController(
	serviceRegistry services.IServiceRegistry,
	sentry sentry.ISentry,
) ISubOrderController {
	return &ISubOrder{
		serviceRegistry: serviceRegistry,
		sentry:          sentry,
	}
}

//nolint:dupl
func (o *ISubOrder) CreateOrder(c *gin.Context) {
	const logCtx = "controllers.http.suborder.sub_order.CreateOrder"
	var (
		ctx     = c.Request.Context()
		request = orderDTO.SubOrderRequest{}
		span    = o.sentry.StartSpan(ctx, logCtx)
	)
	ctx = o.sentry.SpanContext(span)
	defer o.sentry.Finish(span)

	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code:   http.StatusBadRequest,
			Err:    err,
			Gin:    c,
			Sentry: o.sentry,
		})
		return
	}

	validate := validator.New()
	if err = validate.Struct(request); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errorResponse := errorValidation.ErrorValidationResponse(err)
		response.HTTPResponse(response.ParamHTTPResp{
			Err:     err,
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errorResponse,
			Sentry:  o.sentry,
			Gin:     c,
		})
		return
	}

	order, err := o.serviceRegistry.GetSubOrder().CreateOrder(ctx, &request)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code:   http.StatusBadRequest,
			Err:    err,
			Gin:    c,
			Sentry: o.sentry,
		})
		return
	}

	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: order,
		Err:  err,
		Gin:  c,
	})
}

//nolint:dupl
func (o *ISubOrder) GetSubOrderList(c *gin.Context) {
	const logCtx = "controllers.http.suborder.sub_order.GetSubOrderList"
	var (
		ctx     = c.Request.Context()
		request = orderDTO.SubOrderRequestParam{}
		span    = o.sentry.StartSpan(ctx, logCtx)
	)
	ctx = o.sentry.SpanContext(span)
	defer o.sentry.Finish(span)

	err := c.ShouldBindQuery(&request)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code:   http.StatusBadRequest,
			Err:    err,
			Gin:    c,
			Sentry: o.sentry,
		})
		return
	}

	validate := validator.New()
	if err = validate.Struct(request); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errorResponse := errorValidation.ErrorValidationResponse(err)
		response.HTTPResponse(response.ParamHTTPResp{
			Err:     err,
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errorResponse,
			Sentry:  o.sentry,
			Gin:     c,
		})
		return
	}

	order, err := o.serviceRegistry.GetSubOrder().GetSubOrderList(ctx, &request)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code:   http.StatusBadRequest,
			Err:    err,
			Gin:    c,
			Sentry: o.sentry,
		})
		return
	}

	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: order,
		Err:  err,
		Gin:  c,
	})
}

func (o *ISubOrder) GetSubOrderDetail(c *gin.Context) {
	const logCtx = "controllers.http.suborder.sub_order.GetSubOrderDetail"
	var (
		ctx       = c.Request.Context()
		orderUUID = c.Param("uuid")
		span      = o.sentry.StartSpan(ctx, logCtx)
	)
	ctx = o.sentry.SpanContext(span)
	defer o.sentry.Finish(span)

	order, err := o.serviceRegistry.GetSubOrder().GetOrderDetail(ctx, orderUUID)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code:   http.StatusBadRequest,
			Err:    err,
			Gin:    c,
			Sentry: o.sentry,
		})
		return
	}

	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: order,
		Err:  err,
		Gin:  c,
	})
}

func (o *ISubOrder) CancelOrder(c *gin.Context) {
	const logCtx = "controllers.http.suborder.sub_order.CancelOrder"
	var (
		ctx       = c.Request.Context()
		orderUUID = c.Param("uuid")
		span      = o.sentry.StartSpan(ctx, logCtx)
	)
	ctx = o.sentry.SpanContext(span)
	defer o.sentry.Finish(span)

	err := o.serviceRegistry.GetSubOrder().Cancel(ctx, orderUUID)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code:   http.StatusBadRequest,
			Err:    err,
			Gin:    c,
			Sentry: o.sentry,
		})
		return
	}

	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: nil,
		Err:  err,
		Gin:  c,
	})
}
