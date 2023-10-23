package repositories

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
	Create(context.Context, *gorm.DB, *orderPaymentDTO.OrderPaymentRequest) error
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
) error {
	location, _ := time.LoadLocation("Asia/Jakarta") //nolint:errcheck
	datetime := time.Now().In(location)

	orderPayment := orderPaymentModel.OrderPayment{
		SubOrderID: request.SubOrderID,
		PaymentID:  request.PaymentID,
		PaymentURL: request.PaymentLink,
		Status:     request.Status,
		CreatedAt:  &datetime,
		UpdatedAt:  &datetime,
	}
	err := tx.WithContext(ctx).Create(&orderPayment).Error
	if err != nil {
		return errorHelper.WrapError(errorGeneral.ErrSQLError)
	}
	return nil
}
