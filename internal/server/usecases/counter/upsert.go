package counter

import (
	"context"
	"fmt"
	"github.com/aridae/go-metrics-store/internal/server/models"
	"time"
)

func (h *Handler) Upsert(ctx context.Context, updater models.ScalarMetricUpdater) error {
	now := h.now()

	// начать транзакцию

	prevMetricState, err := h.metricsRepo.GetLatestState(ctx, updater.Key())
	if err != nil {
		return fmt.Errorf("metricsRepo.GetScalarMetricLatestState <metricKey:%s>: %w", updater.Key(), err)
	}

	newCounter := buildNewCounter(prevMetricState, updater, now)

	err = h.metricsRepo.Save(ctx, newCounter)
	if err != nil {
		return fmt.Errorf("metricsRepo.Save: %w", err)
	}

	// завершить транзакцию

	return nil
}

func buildNewCounter(
	prevMetricState *models.ScalarMetric,
	metricStateUpdater models.ScalarMetricUpdater,
	now time.Time,
) models.ScalarMetric {
	newState := models.ScalarMetric{
		ScalarMetricUpdater: metricStateUpdater,
		Datetime:            now,
	}

	if prevMetricState != nil {
		newState.Value = prevMetricState.AsCounterValue() + metricStateUpdater.AsCounterValue()
	}

	return newState
}
