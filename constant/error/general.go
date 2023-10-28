package error

import "errors"

const (
	Success = "success"
	Error   = "error"
)

var (
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrSQLError                = errors.New("error sql")
	ErrOrderDate               = errors.New("order date must be greater than now")
	ErrStatus                  = errors.New("invalid status")
)
