package error

import (
	"errors"
)

var (
	ErrOpenState = errors.New("sorry, third party service is busy")
)

var CircuitBreakerErrors = []error{
	ErrOpenState,
}
