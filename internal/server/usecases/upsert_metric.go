package usecases

import (
	"context"
	"fmt"
	"github.com/aridae/go-metrics-store/internal/server/models"
)

func (c *Controller) UpsertMetric(ctx context.Context, metricUpsert models.MetricUpsert) (models.Metric, error) {
	now := c.now()
	var newMetricState models.Metric

	err := c.transactionManager.Do(ctx, func(ctx context.Context) error {
		var txErr error

		newMetricState, txErr = upsert(ctx, c.metricsRepo, metricUpsert, now)
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		return models.Metric{}, fmt.Errorf("transactional upsert metric failed: %w", err)
	}

	return newMetricState, nil
}
