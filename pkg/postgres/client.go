package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

// Client replaceable pgx wrapper
type Client struct {
	connCnf                 *pgxpool.Config
	initialReconnectBackoff time.Duration
	healthcheckTimeout      time.Duration

	poolAcquireTimeout time.Duration
	poolMaxConnections int

	*pgxpool.Pool
}

var defaultOpts = opts{
	initialReconnectBackoff: 15 * time.Second,
	healthcheckTimeout:      5 * time.Second,
	poolAcquireTimeout:      5 * time.Second,
	poolMaxConnections:      15,
}

func NewClient(ctx context.Context, dsn string, opts ...Option) (*Client, error) {
	options := evalOptions(opts...)

	connConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("\tconnCnf pgxpool.ConnConfig\n.ParseDSN: %w", err)
	}

	client := &Client{
		connCnf:                 connConfig,
		initialReconnectBackoff: options.initialReconnectBackoff,
		healthcheckTimeout:      options.healthcheckTimeout,
		poolAcquireTimeout:      options.poolAcquireTimeout,
		poolMaxConnections:      options.poolMaxConnections,
	}

	err = client.connectWithBackoff(ctx, 5)
	if err != nil {
		return nil, err
	}

	return client, nil
}
