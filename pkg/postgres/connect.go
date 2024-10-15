package postgres

import (
	"context"
	"fmt"
	"github.com/aridae/go-metrics-store/internal/server/logger"
	"github.com/jackc/pgx"
	"time"
)

func (c *Client) connectWithBackoff(ctx context.Context, maxRetriesCount int64) error {
	triesLeft := maxRetriesCount
	tryConnectInterval := c.initialReconnectBackoff
	tryConnectAfter := time.NewTimer(c.initialReconnectBackoff)

	for {
		err := c.connect(ctx)
		if err == nil {
			logger.Obtain().Debugf("successfully connected to postgres, happily exiting connectWithBackoff loop")
			return nil
		}

		if triesLeft == 0 {
			return fmt.Errorf("maximum reconnection tries reached, connectWithBackoff loop terminating with error: %w", err)
		}
		triesLeft--

		logger.Obtain().Errorf("error connecting to postgres: %v, will try again after %s", err, tryConnectInterval)

		select {
		case <-ctx.Done():
			logger.Obtain().Infof("stopping connectWithBackoff loop due to context cancel")
			return nil
		case <-tryConnectAfter.C:
			//tryConnectInterval *= 2
			tryConnectAfter.Reset(tryConnectInterval)
		}
	}
}

func (c *Client) connect(_ context.Context) error {
	var err error
	c.ConnPool, err = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     c.connCnf,
		MaxConnections: c.poolMaxConnections,
		AcquireTimeout: c.poolAcquireTimeout,
	})
	if err != nil {
		return fmt.Errorf("could not connect to postgres: %w", err)
	}

	return nil
}
