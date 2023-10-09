package order

import "errors"

var (
	ErrOrderNotFound = errors.New(`error: order not found`)
	ErrOrderNotEmpty = errors.New(`error: previous order not completed yet`)
)
