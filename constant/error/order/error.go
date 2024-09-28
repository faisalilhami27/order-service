package order

import "errors"

var (
	ErrOrderNotFound         = errors.New(`error: order not found`)
	ErrPreviousOrderNotEmpty = errors.New(`error: previous order not completed yet`)
	ErrOrderIsEmpty          = errors.New(`error: order id cannot be empty`)
	ErrCancelOrder           = errors.New(`error: this order already cancelled`)
	ErrInvalidHalfAmount     = errors.New(
		`error: amount must be 50% from (remaining outstanding amount - down payment)`)
	ErrInvalidFullAmount = errors.New(
		`error: amount must be 100% from (remaining outstanding amount - half payment)`)
	ErrFullPaymentNotEmpty = errors.New(`error: your bill for 100% has been paid`)
	ErrHalfPaymentNotEmpty = errors.New(`error: your bill for 50% has been paid`)
)

var OrderErrors = []error{
	ErrOrderNotFound,
	ErrPreviousOrderNotEmpty,
	ErrOrderIsEmpty,
	ErrCancelOrder,
	ErrInvalidHalfAmount,
	ErrInvalidFullAmount,
	ErrFullPaymentNotEmpty,
	ErrHalfPaymentNotEmpty,
}

func CustomErrorOrder(err error) {
	OrderErrors = append(OrderErrors, err)
}
