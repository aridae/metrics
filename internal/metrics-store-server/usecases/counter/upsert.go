package counter

import (
	"context"
	"fmt"
	"github.com/aridae/go-metrics-store/internal/metrics-store-server/models"
	"time"
)

func (h *Handler) Upsert(ctx context.Context, updater models.ScalarMetricUpdater) error {
	now := h.now()

	metricName := updater.Name
	metricType := updater.Type

	// начать транзакцию

	prevMetricState, err := h.metricsRepo.GetLatestState(ctx, metricType, metricName)
	if err != nil {
		return fmt.Errorf("metricsRepo.GetLatestState <metricType:%s> <metricName:%s>: %w", metricType, metricName, err)
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
