package usecases

import (
	"context"
	"fmt"
	"github.com/aridae/go-metrics-store/internal/server/models"
)

func (c *Controller) GetAllScalarMetricsLatestStates(ctx context.Context) ([]models.ScalarMetric, error) {
	metrics, err := c.metricsRepo.GetAllLatestStates(ctx)
	if err != nil {
		return nil, fmt.Errorf("[metricsRepo.GetAllLatestStates] could not get latest metrics states: %w", err)
	}

	return metrics, nil
}
