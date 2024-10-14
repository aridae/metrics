package metricsupsertstrategies

import (
	"context"
	"fmt"
	"time"

	"github.com/aridae/go-metrics-store/internal/server/models"
)

type incrementUpsertStrategy struct{}

func NewIncrementUpsertStrategy() Strategy {
	return &incrementUpsertStrategy{}
}

func (s *incrementUpsertStrategy) Upsert(
	ctx context.Context,
	metricsRepo metricsRepo,
	metricToRegister models.ScalarMetricToRegister,
	now time.Time,
) (models.ScalarMetric, error) {
	// начать транзакцию, которой у меня нет в имитации хранилки в мапке, но будет в постгре

	prevMetricState, err := metricsRepo.GetLatestState(ctx, metricToRegister.Key())
	if err != nil {
		return models.ScalarMetric{}, fmt.Errorf("metricsRepo.GetScalarMetricLatestState <metricKey:%s>: %w", metricToRegister.Key(), err)
	}

	newMetricState, err := increment(prevMetricState, metricToRegister, now)
	if err != nil {
		return models.ScalarMetric{}, fmt.Errorf("increment <metricKey:%s>, failed to build new metric state: %w", metricToRegister.Key(), err)
	}

	err = metricsRepo.Save(ctx, newMetricState)
	if err != nil {
		return models.ScalarMetric{}, fmt.Errorf("metricsRepo.Save <metricKey:%s>: %w", metricToRegister.Key(), err)
	}

	// завершить транзакцию, которой у меня нет в имитации хранилки в мапке, но будет в постгре

	return newMetricState, nil
}

func increment(
	prevMetricState *models.ScalarMetric,
	metricToRegister models.ScalarMetricToRegister,
	now time.Time,
) (models.ScalarMetric, error) {
	if prevMetricState == nil {
		return metricToRegister.AtDatetime(now), nil
	}

	prevValue := prevMetricState.Value()
	incrementValue := metricToRegister.Value()

	newVal, err := prevValue.Inc(incrementValue)
	if err != nil {
		return models.ScalarMetric{}, fmt.Errorf("failed to inc value of prev metric state: %w", err)
	}

	newMetricState := prevMetricState.WithValue(newVal)

	return newMetricState.AtDatetime(now), nil
}
