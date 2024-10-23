package usecases

import (
	"context"
	"fmt"
	"github.com/aridae/go-metrics-store/internal/server/models"
	"github.com/aridae/go-metrics-store/internal/server/repos"
)

func (c *Controller) UpsertMetric(ctx context.Context, metricUpsert models.MetricUpsert) (models.Metric, error) {
	now := c.now()
	var newMetricState models.Metric

	err := c.transactionManager.DoInTransaction(ctx, func(txRepos *repos.Repositories) error {
		var txErr error

		newMetricState, txErr = upsert(ctx, txRepos.MetricRepository, metricUpsert, now)
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
