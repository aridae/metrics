package postgres

import "time"

type opts struct {
	healthcheckTimeout      time.Duration
	initialReconnectBackoff time.Duration
	poolAcquireTimeout      time.Duration
	poolMaxConnections      int
}

type Option func(opts) opts

func WithHealthcheckTimeout(timeout time.Duration) Option {
	return func(o opts) opts {
		o.healthcheckTimeout = timeout
		return o
	}
}

func WithInitialReconnectBackoffOnFail(backoff time.Duration) Option {
	return func(o opts) opts {
		o.initialReconnectBackoff = backoff
		return o
	}
}

func WithMaxConnections(maxConnections int) Option {
	return func(o opts) opts {
		o.poolMaxConnections = maxConnections
		return o
	}
}

func WithAcquireConnectionTimeout(timeout time.Duration) Option {
	return func(o opts) opts {
		o.poolAcquireTimeout = timeout
		return o
	}
}

func evalOptions(options ...Option) opts {
	evalOpts := defaultOpts
	for _, opt := range options {
		evalOpts = opt(evalOpts)
	}

	return evalOpts
}
