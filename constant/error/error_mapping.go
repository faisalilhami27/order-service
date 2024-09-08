package error

import (
	"order-service/constant/error/order"
)

func ErrorMapping(err error) bool {
	allErrors := make([]error, 0)
	allErrors = append(append(GeneralErrors[:], CircuitBreakerErrors[:]...), order.OrderErrors[:]...)

	for _, knownError := range allErrors {
		if err.Error() == knownError.Error() {
			return true
		}
	}
	return false
}
