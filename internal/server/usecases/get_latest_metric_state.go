package usecases

import (
	"context"
	"fmt"
	"github.com/aridae/go-metrics-store/internal/server/models"
)

func (c *Controller) GetScalarMetricLatestState(ctx context.Context, metricKey models.MetricKey) (*models.ScalarMetric, error) {
	metric, err := c.metricsRepo.GetLatestState(ctx, metricKey)
	if err != nil {
		return nil, fmt.Errorf("[metricsRepo.GetScalarMetricLatestState] could not get latest metrics state <key:%s>: %w", metricKey, err)
	}

	return metric, nil
}
