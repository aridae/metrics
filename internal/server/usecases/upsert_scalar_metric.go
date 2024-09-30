package usecases

import (
	"context"
	"fmt"
	metricsupsertstrategies "github.com/aridae/go-metrics-store/internal/server/metrics-upsert-strategies"
	"github.com/aridae/go-metrics-store/internal/server/models"
)

func (c *Controller) UpsertScalarMetric(ctx context.Context, metricToRegister models.ScalarMetricToRegister, strategy metricsupsertstrategies.Strategy) error {
	now := c.now()

	err := strategy.Upsert(ctx, c.metricsRepo, metricToRegister, now)
	if err != nil {
		return fmt.Errorf("c.rtMetricsUpsertStrategy.Upsert: %w", err)
	}

	return nil
}
