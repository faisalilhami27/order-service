package error

import (
	"order-service/constant/error/order"
)

//nolint:revive,ineffassign
func ErrorMapping(err error) bool {
	allErrors := make([]error, 0)                                                                     //nolint:staticcheck
	allErrors = append(append(GeneralErrors[:], CircuitBreakerErrors[:]...), order.OrderErrors[:]...) //nolint:gocritic

	for _, knownError := range allErrors {
		if err.Error() == knownError.Error() {
			return true
		}
	}
	return false
}
