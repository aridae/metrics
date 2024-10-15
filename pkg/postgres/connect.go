package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/aridae/go-metrics-store/internal/server/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func connectWithBackoff(
	ctx context.Context,
	cnf *pgxpool.Config,
	maxRetriesCount int64,
	initialReconnectBackoff time.Duration,
) (*pgxpool.Pool, error) {
	triesLeft := maxRetriesCount
	tryConnectInterval := initialReconnectBackoff
	tryConnectAfter := time.NewTimer(initialReconnectBackoff)

	for {
		pool, err := connect(ctx, cnf)
		if err == nil {
			logger.Obtain().Debugf("successfully connected to postgres, happily exiting connectWithBackoff loop")
			return pool, nil
		}

		if triesLeft == 0 {
			return nil, fmt.Errorf("maximum reconnection tries reached, connectWithBackoff loop terminating with error: %w", err)
		}
		triesLeft--

		logger.Obtain().Errorf("error connecting to postgres: %v, will try again after %s", err, tryConnectInterval)

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("terminating connectWithBackoff loop due to context cancel: %w", ctx.Err())
		case <-tryConnectAfter.C:
			tryConnectInterval *= 2
			tryConnectAfter.Reset(tryConnectInterval)
		}
	}
}

func connect(ctx context.Context, config *pgxpool.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("could not connect to postgres pgx pool: %w", err)
	}

	return pool, nil
}
