package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx"
)

// Client replaceable pgx wrapper
type Client struct {
	connCnf pgx.ConnConfig

	initialReconnectBackoff time.Duration
	healthcheckTimeout      time.Duration

	poolAcquireTimeout time.Duration
	poolMaxConnections int

	*pgx.ConnPool
}

var defaultOpts = opts{
	initialReconnectBackoff: 15 * time.Second,
	healthcheckTimeout:      5 * time.Second,
	poolAcquireTimeout:      5 * time.Second,
	poolMaxConnections:      15,
}

func NewClient(ctx context.Context, dsn string, opts ...Option) (*Client, error) {
	options := evalOptions(opts...)

	connConfig, err := pgx.ParseDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("pgx.ParseDSN: %w", err)
	}

	client := &Client{
		connCnf:                 connConfig,
		initialReconnectBackoff: options.initialReconnectBackoff,
		healthcheckTimeout:      options.healthcheckTimeout,
		poolAcquireTimeout:      options.poolAcquireTimeout,
		poolMaxConnections:      options.poolMaxConnections,
	}

	go client.connectWithBackoff(ctx, 90)

	return client, nil
}
