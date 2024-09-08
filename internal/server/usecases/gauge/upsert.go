package gauge

import (
	"context"
	"fmt"
	"github.com/aridae/go-metrics-store/internal/server/models"
)

func (h *Handler) Upsert(ctx context.Context, updater models.ScalarMetricUpdater) error {
	now := h.now()

	newMetricState := models.ScalarMetric{
		ScalarMetricUpdater: updater,
		Datetime:            now,
	}

	err := h.repo.Save(ctx, newMetricState)
	if err != nil {
		return fmt.Errorf("metricsRepo.Save: %w", err)
	}

	return nil
}
