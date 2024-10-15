package postgres

import (
	"context"
	"fmt"
)

func (c *Client) Healthcheck(ctx context.Context) error {
	if c == nil || c.ConnPool == nil {
		return fmt.Errorf("nil receiver, connection is not established")
	}

	ctx, cancel := context.WithTimeout(ctx, c.healthcheckTimeout)
	defer cancel()

	_, err := c.ConnPool.ExecEx(ctx, ";", nil)
	if err != nil {
		return fmt.Errorf("postgres Connection Pool seems to be unreachable, ExecEx: %w", err)
	}

	return nil
}
