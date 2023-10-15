package error

import "errors"

const (
	Success = "success"
	Error   = "error"
)

var (
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrSQLError                = errors.New("error sql")
)
