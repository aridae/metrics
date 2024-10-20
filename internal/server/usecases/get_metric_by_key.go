package usecases

import (
	"context"
	"fmt"

	"github.com/aridae/go-metrics-store/internal/server/models"
)

func (c *Controller) GetMetricByKey(ctx context.Context, metricKey models.MetricKey) (*models.Metric, error) {
	metric, err := c.metricsRepo.GetByKey(ctx, metricKey)
	if err != nil {
		return nil, fmt.Errorf("[metricsRepo.GetByKey] failed to get latest metric state <key:%s>: %w", metricKey, err)
	}

	return metric, nil
}
