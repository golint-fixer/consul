package consul

import (
	"time"

	retry "gopkg.in/eapache/go-resiliency.v1/retrier"
)

const (
	// RetryTimes defines the default max amount of times to retry a request.
	RetryTimes = 5

	// RetryWait defines the default amount of time to wait before each retry attempt.
	RetryWait = 100 * time.Millisecond
)

var (
	// ConstantBackoff provides a built-in retry strategy based on constant back off.
	ConstantBackoff = retry.New(retry.ConstantBackoff(RetryTimes, RetryWait), nil)

	// ExponentialBackoff provides a built-int retry strategy based on exponential back off.
	ExponentialBackoff = retry.New(retry.ExponentialBackoff(RetryTimes, RetryWait), nil)

	// DefaultRetrier stores the default retry strategy used by the plugin.
	// By default will use a constant retry strategy with a maximum of 3 retry attempts.
	DefaultRetrier = ConstantBackoff
)

// Retrier defines the required interface implemented by retry strategies.
type Retrier interface {
	Run(func() error) error
}

// Retry provides a retry.Retrier capable interface that
// encapsulates Consul client and user defined strategy.
type Retry struct {
	// retrier stores the retry strategy to be used.
	retrier Retrier
}

// NewRetrier creates a default retrier for the given Consul client and context.
func NewRetrier(r Retrier) *Retry {
	return &Retry{retrier: r}
}

// Run runs the given function multiple times, acting like a proxy
// to user defined retry strategy.
func (r *Retry) Run(fn func() error) error {
	return r.retrier.Run(fn)
}
