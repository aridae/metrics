package usecases

import (
	"context"
	"github.com/aridae/go-metrics-store/internal/server/models"
)

func (c *Controller) UpsertMetric(ctx context.Context, metricUpsert models.MetricUpsert) (models.Metric, error) {
	now := c.now()

	newMetricState, err := c.upsert(ctx, metricUpsert, now)
	if err != nil {
		return models.Metric{}, err
	}

	return newMetricState, nil
}
