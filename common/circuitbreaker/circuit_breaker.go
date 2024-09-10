package circuitbreaker

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/sony/gobreaker"

	"order-service/common/sentry"
	errCircuitBreaker "order-service/constant/error"
)

type CircuitBreaker struct {
	name        string
	maxRequests uint32
	timeout     uint32
	cb          *gobreaker.CircuitBreaker
	sentry      sentry.ISentry
}

type BreakerFunc func() (interface{}, error)
type Option func(*CircuitBreaker)

type ICircuitBreaker interface {
	Execute(context.Context, BreakerFunc) error
}

func WithMaxRequest(maxRequest uint32) Option {
	return func(c *CircuitBreaker) {
		c.maxRequests = maxRequest
	}
}

func WithTimeout(timeout uint32) Option {
	return func(c *CircuitBreaker) {
		c.timeout = timeout
	}
}

func NewCircuitBreaker(sentry sentry.ISentry, options ...Option) ICircuitBreaker {
	circuitBreaker := &CircuitBreaker{
		name:   "circuit-breaker",
		sentry: sentry,
	}

	for _, option := range options {
		option(circuitBreaker)
	}

	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        circuitBreaker.name,
		MaxRequests: circuitBreaker.maxRequests,
		Interval:    time.Duration(circuitBreaker.timeout) * time.Second,
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) { //nolint:revive
			log.Infof("Circuit Breaker state changed from %s to %s\n", from, to)
		},
	})

	circuitBreaker.cb = cb
	return circuitBreaker
}

func (c CircuitBreaker) Execute(ctx context.Context, client BreakerFunc) error {
	logCtx := "common.circuitbreaker.circuit_breaker.Execute"
	var (
		span = c.sentry.StartSpan(ctx, logCtx)
	)
	c.sentry.SpanContext(span)
	defer c.sentry.Finish(span)

	result, err := c.cb.Execute(client)
	if err != nil {
		if c.cb.State() == gobreaker.StateOpen {
			return errCircuitBreaker.ErrOpenState
		}

		return err
	}

	log.Infof("Result: %v\n", result)
	return nil
}
