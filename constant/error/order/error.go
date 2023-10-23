package order

import "errors"

var (
	ErrOrderNotFound = errors.New(`error: suborder not found`)
	ErrOrderNotEmpty = errors.New(`error: previous suborder not completed yet`)
	ErrCancelOrder   = errors.New(`error: this suborder already cancelled`)
)
