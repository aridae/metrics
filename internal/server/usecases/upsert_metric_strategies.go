package usecases

import (
	"context"
	"fmt"
	"github.com/aridae/go-metrics-store/internal/server/models"
	"time"
)

func upsert(ctx context.Context, metricsRepo metricsRepo, metricUpsert models.MetricUpsert, now time.Time) (models.Metric, error) {
	upsertFn, ok := upsertStrategyByType[metricUpsert.GetType()]
	if !ok {
		return models.Metric{}, fmt.Errorf("unknown metric type: %s", metricUpsert.GetType())
	}

	newMetricState, err := upsertFn(ctx, metricsRepo, metricUpsert, now)
	if err != nil {
		return models.Metric{}, fmt.Errorf("failed to upsert metric <key:%s>: %w", metricUpsert.GetKey(), err)
	}

	return newMetricState, nil
}

var upsertStrategyByType = map[models.MetricType]func(ctx context.Context, metricsRepo metricsRepo, metricUpsert models.MetricUpsert, now time.Time) (models.Metric, error){
	models.MetricTypeGauge:   upsertMetricOverride,
	models.MetricTypeCounter: upsertMetricIncrement,
}

func upsertMetricOverride(ctx context.Context, metricsRepo metricsRepo, metricUpsert models.MetricUpsert, now time.Time) (models.Metric, error) {
	newState := metricUpsert.WithDatetime(now)

	err := metricsRepo.Save(ctx, newState)
	if err != nil {
		return models.Metric{}, fmt.Errorf("failed to save new metric state <key:%s>: %w", metricUpsert.GetKey(), err)
	}

	return newState, nil
}

func upsertMetricIncrement(ctx context.Context, metricsRepo metricsRepo, metricUpsert models.MetricUpsert, now time.Time) (models.Metric, error) {
	prevState, err := metricsRepo.GetByKey(ctx, metricUpsert.GetKey())
	if err != nil {
		return models.Metric{}, fmt.Errorf("failed to get prev metric state <key:%s>: %w", metricUpsert.GetKey(), err)
	}

	newState := metricUpsert.WithDatetime(now)
	if prevState != nil {
		newVal, err := prevState.GetValue().Inc(metricUpsert.GetValue())
		if err != nil {
			return models.Metric{}, fmt.Errorf("failed to do increment on prev metric state <key:%s>: %w", metricUpsert.GetKey(), err)
		}

		newState = prevState.WithValue(newVal).WithDatetime(now)
	}

	err = metricsRepo.Save(ctx, newState)
	if err != nil {
		return models.Metric{}, fmt.Errorf("failed to save new metric state <key:%s>: %w", metricUpsert.GetKey(), err)
	}

	return newState, nil
}
