package usecases

import (
	"context"
	"fmt"

	"github.com/aridae/go-metrics-store/internal/server/models"
)

// GetAllMetrics получает все метрики из репозитория.
//
// Аргументы:
// ctx (context.Context): Контекст выполнения запроса.
//
// Возвращает:
// []models.Metric: Список всех метрик.
// error: Ошибка, если что-то пошло не так при получении метрик.
func (c *Controller) GetAllMetrics(ctx context.Context) ([]models.Metric, error) {
	metrics, err := c.metricsRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("[metricsRepo.GetAll] failed to get latest metrics-reporting states: %w", err)
	}

	return metrics, nil
}
