package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/aridae/go-metrics-store/internal/server/logger"
	"github.com/jackc/pgx"
)

// Client replaceable pgx wrapper
type Client struct {
	*pgx.ConnPool
	healthCheckTimeout time.Duration
}

func NewClient(dsn string, maxOpenConn int) (*Client, error) {
	connConfig, err := pgx.ParseDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("pgx.ParseDSN: %w", err)
	}

	defaultTimeout := time.Second * 3

	conn, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     connConfig,
		MaxConnections: maxOpenConn,
		AfterConnect: func(conn *pgx.Conn) error {
			logger.Obtain().Debugf("got new connection pid=%d", conn.PID())
			return nil
		},
		AcquireTimeout: defaultTimeout,
	})
	if err != nil {
		return nil, fmt.Errorf("pgx.NewConnPool: %w", err)
	}

	return &Client{ConnPool: conn, healthCheckTimeout: defaultTimeout}, nil
}

func (c *Client) Healthcheck(ctx context.Context) error {
	if c == nil {
		return fmt.Errorf("nil client receiver")
	}

	ctx, cancel := context.WithTimeout(ctx, c.healthCheckTimeout)
	defer cancel()

	_, err := c.ConnPool.ExecEx(ctx, ";", nil)
	if err != nil {
		return fmt.Errorf("postgres Connection Pool seems to be unreachable, ExecEx: %w", err)
	}

	return nil
}
