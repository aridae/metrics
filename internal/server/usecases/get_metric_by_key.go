package usecases

import (
	"context"
	"fmt"

	"github.com/aridae/go-metrics-store/internal/server/models"
)

// GetMetricByKey получает метрику по заданному ключу из репозитория.
//
// Аргументы:
// ctx (context.Context): Контекст выполнения запроса.
// metricKey (models.MetricKey): Ключ метрики, которую нужно получить.
//
// Возвращает:
// *models.Metric: Метрика с указанным ключом.
// error: Ошибка, если что-то пошло не так при получении метрики.
func (c *Controller) GetMetricByKey(ctx context.Context, metricKey models.MetricKey) (*models.Metric, error) {
	metric, err := c.metricsRepo.GetByKey(ctx, metricKey)
	if err != nil {
		return nil, fmt.Errorf("[metricsRepo.GetByKey] failed to get latest metric state <key:%s>: %w", metricKey, err)
	}

	return metric, nil
}
