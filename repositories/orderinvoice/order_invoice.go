package repositories

import (
	"context"

	"order-service/common/sentry"
	orderInvoiceModel "order-service/domain/models"

	"time"

	"gorm.io/gorm"

	errorGeneral "order-service/constant/error"
	errorHelper "order-service/utils/error"
)

type IOrderInvoice struct {
	db     *gorm.DB
	sentry sentry.ISentry
}

type IOrderInvoiceRepository interface {
	Create(context.Context, *gorm.DB, *orderInvoiceModel.OrderInvoice) error
}

func NewOrderInvoice(db *gorm.DB, sentry sentry.ISentry) IOrderInvoiceRepository {
	return &IOrderInvoice{
		db:     db,
		sentry: sentry,
	}
}

func (o *IOrderInvoice) Create(
	ctx context.Context,
	tx *gorm.DB,
	request *orderInvoiceModel.OrderInvoice,
) error {
	const logCtx = "repositories.orderinvoice.order_invoice.Create"
	var (
		span         = o.sentry.StartSpan(ctx, logCtx)
		orderInvoice orderInvoiceModel.OrderInvoice
	)
	ctx = o.sentry.SpanContext(span)
	defer o.sentry.Finish(span)

	location, _ := time.LoadLocation("Asia/Jakarta") //nolint:errcheck
	datetime := time.Now().In(location)

	orderInvoice = orderInvoiceModel.OrderInvoice{
		SubOrderID:    request.SubOrderID,
		InvoiceID:     request.InvoiceID,
		InvoiceNumber: request.InvoiceNumber,
		InvoiceURL:    request.InvoiceURL,
		CreatedAt:     &datetime,
		UpdatedAt:     &datetime,
	}
	err := tx.WithContext(ctx).Create(&orderInvoice).Error
	if err != nil {
		return errorHelper.WrapError(errorGeneral.ErrSQLError, o.sentry)
	}
	return nil
}
