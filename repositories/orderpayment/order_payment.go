package repositories

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"order-service/common/sentry"
	orderPaymentModel "order-service/domain/models"

	"time"

	errorGeneral "order-service/constant/error"
	orderPaymentDTO "order-service/domain/dto/orderpayment"
	errorHelper "order-service/utils/error"
)

type IOrderPayment struct {
	db     *gorm.DB
	sentry sentry.ISentry
}

type IOrderPaymentRepository interface {
	Create(context.Context, *gorm.DB, *orderPaymentDTO.OrderPaymentRequest) error
	Update(context.Context, *gorm.DB, *orderPaymentDTO.OrderPaymentRequest) error
	FindByPaymentID(context.Context, *gorm.DB, string) (*orderPaymentModel.OrderPayment, error)
}

func NewOrderPayment(db *gorm.DB, sentry sentry.ISentry) IOrderPaymentRepository {
	return &IOrderPayment{
		db:     db,
		sentry: sentry,
	}
}

func (o *IOrderPayment) Create(
	ctx context.Context,
	tx *gorm.DB,
	request *orderPaymentDTO.OrderPaymentRequest,
) error {
	const logCtx = "repositories.orderpayment.order_payment.Create"
	var (
		span         = o.sentry.StartSpan(ctx, logCtx)
		orderPayment orderPaymentModel.OrderPayment
	)
	ctx = o.sentry.SpanContext(span)
	defer o.sentry.Finish(span)

	location, _ := time.LoadLocation("Asia/Jakarta") //nolint:errcheck
	datetime := time.Now().In(location)

	orderPayment = orderPaymentModel.OrderPayment{
		Amount:      request.Amount,
		SubOrderID:  request.SubOrderID,
		PaymentID:   request.PaymentID,
		PaymentURL:  &request.PaymentLink,
		PaymentType: request.PaymentType,
		VANumber:    request.VANumber,
		Bank:        request.Bank,
		Acquirer:    request.Acquirer,
		Status:      request.Status,
		ExpiredAt:   request.ExpiredAt,
		PaidAt:      request.PaidAt,
		CreatedAt:   &datetime,
		UpdatedAt:   &datetime,
	}
	err := tx.WithContext(ctx).Create(&orderPayment).Error
	if err != nil {
		return errorHelper.WrapError(errorGeneral.ErrSQLError, o.sentry)
	}
	return nil
}

func (o *IOrderPayment) Update(
	ctx context.Context,
	tx *gorm.DB,
	request *orderPaymentDTO.OrderPaymentRequest,
) error {
	const logCtx = "repositories.orderpayment.order_payment.Update"
	var (
		span         = o.sentry.StartSpan(ctx, logCtx)
		orderPayment orderPaymentModel.OrderPayment
	)
	ctx = o.sentry.SpanContext(span)
	defer o.sentry.Finish(span)

	orderPayment = orderPaymentModel.OrderPayment{
		Amount:      request.Amount,
		PaymentID:   request.PaymentID,
		PaymentURL:  &request.PaymentLink,
		PaymentType: request.PaymentType,
		VANumber:    request.VANumber,
		Bank:        request.Bank,
		Acquirer:    request.Acquirer,
		Status:      request.Status,
		PaidAt:      request.PaidAt,
	}
	err := tx.WithContext(ctx).
		Where("payment_id = ?", request.PaymentID).
		Updates(&orderPayment).Error
	if err != nil {
		return errorHelper.WrapError(errorGeneral.ErrSQLError, o.sentry)
	}
	return nil
}

func (o *IOrderPayment) FindByPaymentID(
	ctx context.Context,
	tx *gorm.DB,
	paymentID string,
) (*orderPaymentModel.OrderPayment, error) {
	const logCtx = "repositories.orderpayment.order_payment.Update"
	var (
		span         = o.sentry.StartSpan(ctx, logCtx)
		orderPayment orderPaymentModel.OrderPayment
	)
	ctx = o.sentry.SpanContext(span)
	defer o.sentry.Finish(span)

	err := tx.WithContext(ctx).
		Where("payment_id = ?", paymentID).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&orderPayment).Error
	if err != nil {
		return nil, errorHelper.WrapError(errorGeneral.ErrSQLError, o.sentry)
	}
	return &orderPayment, nil
}
