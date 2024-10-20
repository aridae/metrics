package postgres

import "time"

type opts struct {
	initialReconnectBackoff time.Duration
}

type Option func(opts) opts

func WithInitialReconnectBackoffOnFail(backoff time.Duration) Option {
	return func(o opts) opts {
		o.initialReconnectBackoff = backoff
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
