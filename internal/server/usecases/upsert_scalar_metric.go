package usecases

import (
	"context"
	"fmt"
	"github.com/aridae/go-metrics-store/internal/server/models"
)

func (c *Controller) UpsertScalarMetric(ctx context.Context, updater models.ScalarMetricUpdater) error {
	switch updater.Type {
	case models.ScalarMetricTypeCounter:
		return c.counterUseCasesHandler.Upsert(ctx, updater)
	case models.ScalarMetricTypeGauge:
		return c.gaugeUseCasesHandler.Upsert(ctx, updater)
	default:
		return fmt.Errorf("unknown scalar-metrics type: %v", updater.Type)
	}
}
