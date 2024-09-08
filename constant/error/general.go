package error

import "errors"

const (
	Success = "success"
	Error   = "error"
)

var (
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrSQLError                = errors.New("database server failed to execute, please try again")
	ErrOrderDate               = errors.New("order date must be greater than now")
	ErrStatus                  = errors.New("invalid status")
	ErrTooManyRequest          = errors.New("too many request, please try again later")
	ErrUnauthorized            = errors.New("unauthorized")
	ErrForbidden               = errors.New("you don't have permission to access this resource")
)

var GeneralErrors = []error{
	ErrInvalidStatusTransition,
	ErrSQLError,
	ErrOrderDate,
	ErrStatus,
	ErrTooManyRequest,
	ErrUnauthorized,
	ErrForbidden,
}
