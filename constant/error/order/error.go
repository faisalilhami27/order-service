package order

import "errors"

var (
	ErrOrderNotFound         = errors.New(`error: suborder not found`)
	ErrPreviousOrderNotEmpty = errors.New(`error: previous order not completed yet`)
	ErrOrderIsEmpty          = errors.New(`error: order id cannot be empty`)
	ErrCancelOrder           = errors.New(`error: this order already cancelled`)
	ErrHalfPaymentIsEmpty    = errors.New(`error: you must be pay half payment first`)
	ErrInvalidDownAmount     = errors.New(`error: amount must be 10% from total price`)
	ErrInvalidHalfAmount     = errors.New(
		`error: amount must be 50% from (remaining outstanding amount - down payment)`)
	ErrInvalidFullAmount = errors.New(
		`error: amount must be 100% from (remaining outstanding amount - half payment)`)
)
