package usecases

import (
	"context"
	"github.com/aridae/go-metrics-store/internal/server/models"
)

func (c *Controller) UpsertMetricsBatch(ctx context.Context, metricsUpserts []models.MetricUpsert) ([]models.Metric, error) {
	now := c.now()

	newMetricStates := make([]models.Metric, 0, len(metricsUpserts))
	for _, metricUpsert := range metricsUpserts {
		newMetricState, err := c.upsert(ctx, metricUpsert, now)
		if err != nil {
			return nil, err
		}

		newMetricStates = append(newMetricStates, newMetricState)
	}

	return newMetricStates, nil
}
