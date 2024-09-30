package metricsupsertstrategies

import (
	"context"
	"fmt"
	"github.com/aridae/go-metrics-store/internal/server/models"
	"time"
)

type overrideUpsertStrategy struct{}

func NewOverrideUpsertStrategy() Strategy {
	return &overrideUpsertStrategy{}
}

func (s *overrideUpsertStrategy) Upsert(ctx context.Context, metricsRepo metricsRepo, metricToRegister models.ScalarMetricToRegister, now time.Time) error {
	newMetricState := metricToRegister.AtDatetime(now)

	err := metricsRepo.Save(ctx, newMetricState)
	if err != nil {
		return fmt.Errorf("metricsRepo.Save: %w", err)
	}

	return nil
}
