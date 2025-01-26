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

// upsertMetricOverride обновляет состояние метрики или создает новое, если оно отсутствует.
//
// Аргументы:
// ctx (context.Context): Контекст выполнения запроса.
// metricsRepo (metricsRepo): Репозиторий для работы с метриками.
// metricUpsert (models.MetricUpsert): Данные для обновления или создания метрики.
// now (time.Time): Текущее время.
//
// Возвращает:
// models.Metric: Обновленное или созданное состояние метрики.
// error: Ошибка, если что-то пошло не так при сохранении состояния метрики.
func upsertMetricOverride(ctx context.Context, metricsRepo metricsRepo, metricUpsert models.MetricUpsert, now time.Time) (models.Metric, error) {
	newState := metricUpsert.WithDatetime(now)

	err := metricsRepo.Save(ctx, newState)
	if err != nil {
		return models.Metric{}, fmt.Errorf("failed to save new metric state <key:%s>: %w", metricUpsert.GetKey(), err)
	}

	return newState, nil
}

// upsertMetricIncrement увеличивает значение метрики или создает новую, если она отсутствует.
//
// Аргументы:
// ctx (context.Context): Контекст выполнения запроса.
// metricsRepo (metricsRepo): Репозиторий для работы с метриками.
// metricUpsert (models.MetricUpsert): Данные для увеличения или создания метрики.
// now (time.Time): Текущее время.
//
// Возвращает:
// models.Metric: Увеличившая или созданная метрика.
// error: Ошибка, если что-то пошло не так при увеличении или сохранении метрики.
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
