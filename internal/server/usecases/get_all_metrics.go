package usecases

import (
	"context"
	"fmt"

	"github.com/aridae/go-metrics-store/internal/server/models"
)

func (c *Controller) GetAllMetrics(ctx context.Context) ([]models.Metric, error) {
	metrics, err := c.metricsRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("[metricsRepo.GetAll] failed to get latest metrics-reporting states: %w", err)
	}

	return metrics, nil
}
