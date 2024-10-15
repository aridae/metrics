package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

// Client replaceable pgx wrapper
type Client struct {
	healthcheckTimeout time.Duration
	*pgxpool.Pool
}

var defaultOpts = opts{
	initialReconnectBackoff: 1 * time.Second,
	healthcheckTimeout:      5 * time.Second,
}

func NewClient(ctx context.Context, dsn string, opts ...Option) (*Client, error) {
	options := evalOptions(opts...)

	connConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}

	pool, err := connectWithBackoff(ctx, connConfig, 5, options.initialReconnectBackoff)
	if err != nil {
		return nil, fmt.Errorf("connectWithBackoff: %w", err)
	}

	return &Client{
		healthcheckTimeout: options.healthcheckTimeout,
		Pool:               pool,
	}, nil
}
