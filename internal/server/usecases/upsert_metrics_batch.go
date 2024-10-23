package usecases

import (
	"context"
	"fmt"
	"github.com/aridae/go-metrics-store/internal/server/models"
	"github.com/aridae/go-metrics-store/internal/server/repos"
)

func (c *Controller) UpsertMetricsBatch(ctx context.Context, metricsUpserts []models.MetricUpsert) ([]models.Metric, error) {
	now := c.now()
	newMetricStates := make([]models.Metric, 0, len(metricsUpserts))

	err := c.transactionManager.DoInTransaction(ctx, func(txRepos *repos.Repositories) error {
		for _, metricUpsert := range metricsUpserts {

			newMetricState, txErr := upsert(ctx, txRepos.MetricRepository, metricUpsert, now)
			if txErr != nil {
				return txErr
			}

			newMetricStates = append(newMetricStates, newMetricState)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("transactional upsert metric failed: %w", err)
	}

	return newMetricStates, nil
}
