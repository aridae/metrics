package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/aridae/go-metrics-store/internal/server/models"
	metricinmemrepo "github.com/aridae/go-metrics-store/internal/server/repos/metric/metric-inmem-repo"
	"github.com/aridae/go-metrics-store/pkg/inmem"
	nooptrm "github.com/aridae/go-metrics-store/pkg/noop-trm"
)

// ExampleController_GetAllMetrics демонстрирует использование метода Controller.GetAllMetrics.
func ExampleController_GetAllMetrics() {
	ctx := context.Background()

	metricRepo := metricinmemrepo.NewRepositoryImplementation(inmem.New[models.MetricKey, models.Metric]())

	err := metricRepo.Save(ctx, models.Metric{
		MetricUpsert: models.MetricUpsert{
			MName: "name-1",
			Mtype: models.MetricTypeCounter,
			Val:   models.NewInt64MetricValue(123),
		},
	})
	if err != nil {
		fmt.Printf("Ошибка сохранения метрики %s: %v\n", "name-1", err)
		return
	}

	err = metricRepo.Save(ctx, models.Metric{
		MetricUpsert: models.MetricUpsert{
			MName: "name-2",
			Mtype: models.MetricTypeGauge,
			Val:   models.NewFloat64MetricValue(123.5),
		},
	})
	if err != nil {
		fmt.Printf("Ошибка сохранения метрики %s: %v\n", "name-2", err)
		return
	}

	// Создайте экземпляр Controller с настроенным репозиторием
	c := &Controller{
		metricsRepo:        metricRepo, // Подставьте реальный репозиторий
		transactionManager: nooptrm.NewNoopTransactionManager(),
		now:                time.Now().UTC,
	}

	// Вызовите метод GetAllMetrics
	metrics, err := c.GetAllMetrics(ctx)
	if err != nil {
		fmt.Printf("Ошибка получения метрик: %v\n", err)
		return
	}

	// Выведем полученные метрики
	for _, metric := range metrics {
		fmt.Printf("Метрика Key: %s Val: %s\n", metric.GetKey(), metric.GetValue().String())
	}
	// Output:
	// Метрика Key: counter:name-1 Val: 123
	// Метрика Key: gauge:name-2 Val: 123.5
}
