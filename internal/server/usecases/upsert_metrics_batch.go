package usecases

import (
	"context"
	"fmt"

	"github.com/aridae/go-metrics-store/internal/server/models"
)

// UpsertMetricsBatch обновляет несколько существующих метрик или создает новые, если они отсутствуют.
//
// Аргументы:
// ctx (context.Context): Контекст выполнения запроса.
// metricsUpserts ([]models.MetricUpsert): Массив данных для обновления или создания метрик.
//
// Возвращает:
// []models.Metric: Массив обновленных или новых метрик.
// error: Ошибка, если что-то пошло не так во время транзакции.
func (c *Controller) UpsertMetricsBatch(ctx context.Context, metricsUpserts []models.MetricUpsert) ([]models.Metric, error) {
	now := c.now()
	newMetricStates := make([]models.Metric, 0, len(metricsUpserts))

	err := c.transactionManager.Do(ctx, func(ctx context.Context) error {
		for _, metricUpsert := range metricsUpserts {

			newMetricState, txErr := upsert(ctx, c.metricsRepo, metricUpsert, now)
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
