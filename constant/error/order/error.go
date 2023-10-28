package order

import "errors"

var (
	ErrOrderNotFound = errors.New(`error: suborder not found`)
	ErrOrderNotEmpty = errors.New(`error: previous order not completed yet`)
	ErrCancelOrder   = errors.New(`error: this order already cancelled`)
)
