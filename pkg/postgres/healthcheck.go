package postgres

import (
	"context"
	"fmt"
)

func (c *Client) Healthcheck(ctx context.Context) error {
	if c == nil || c.Pool == nil {
		return fmt.Errorf("nil receiver, connection is not established")
	}

	ctx, cancel := context.WithTimeout(ctx, c.healthcheckTimeout)
	defer cancel()

	err := c.Pool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("postgres Connection Pool seems to be unreachable, Ping failed: %w", err)
	}

	return nil
}
