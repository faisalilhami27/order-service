package orderpayment

import (
	"context"
	"time"

	"gorm.io/gorm"

	errorGeneral "order-service/constant/error"
	orderPaymentDTO "order-service/domain/dto/orderpayment"
	orderPaymentModel "order-service/domain/models/orderpayment"
	errorHelper "order-service/utils/error"
)

type IOrderPayment struct {
	db *gorm.DB
}

type IOrderPaymentRepository interface {
	Create(context.Context, *gorm.DB, *orderPaymentDTO.OrderPaymentRequest) (*orderPaymentModel.OrderPayment, error)
}

func NewOrderPayment(db *gorm.DB) IOrderPaymentRepository {
	return &IOrderPayment{
		db: db,
	}
}

func (o *IOrderPayment) Create(
	ctx context.Context,
	tx *gorm.DB,
	request *orderPaymentDTO.OrderPaymentRequest,
) (*orderPaymentModel.OrderPayment, error) {
	location, _ := time.LoadLocation("Asia/Jakarta") //nolint:errcheck
	datetime := time.Now().In(location)

	orderPayment := orderPaymentModel.OrderPayment{
		OrderID:    request.OrderID,
		PaymentID:  request.PaymentID,
		InvoiceID:  request.InvoiceID,
		PaymentURL: request.PaymentURL,
		Status:     request.Status,
		CreatedAt:  &datetime,
		UpdatedAt:  &datetime,
	}
	err := tx.WithContext(ctx).Create(&orderPayment).Error
	if err != nil {
		return nil, errorHelper.WrapError(errorGeneral.ErrSQLError)
	}
	return &orderPayment, nil
}
