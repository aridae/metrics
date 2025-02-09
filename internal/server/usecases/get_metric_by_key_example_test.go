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

// ExampleController_GetMetricByKey демонстрирует использование метода Controller.GetMetricByKey.
func ExampleController_GetMetricByKey() {
	ctx := context.Background()

	exampleMetric := models.Metric{
		MetricUpsert: models.MetricUpsert{
			MName: "name-1",
			Mtype: models.MetricTypeCounter,
			Val:   models.NewInt64MetricValue(123),
		},
	}
	exampleMetricKey := exampleMetric.GetKey()

	metricRepo := metricinmemrepo.NewRepositoryImplementation(inmem.New[models.MetricKey, models.Metric]())

	err := metricRepo.Save(ctx, exampleMetric)
	if err != nil {
		fmt.Printf("Ошибка сохранения метрики %s: %v\n", "name-1", err)
		return
	}

	// Создайте экземпляр Controller с настроенным репозиторием
	c := &Controller{
		metricsRepo:        metricRepo, // Подставьте реальный репозиторий
		transactionManager: nooptrm.NewNoopTransactionManager(),
		now:                time.Now().UTC,
	}

	// Вызовите метод GetMetricByKey, ожидаем получить метрику по ключу
	foundMetric, err := c.GetMetricByKey(ctx, exampleMetricKey)
	if err != nil {
		fmt.Printf("Ошибка получения метрики %s: %v\n", exampleMetricKey, err)
		return
	}
	fmt.Printf("Найдена метрика Key: %s Val: %s\n", exampleMetricKey, foundMetric.GetValue().String())

	// Вызовите метод GetMetricByKey, ожидаем получить nil в случае, если метрика не найдена
	exampleNotFoundMetricKey := models.BuildMetricKey("not-found-metric", models.MetricTypeGauge)
	notFoundMetric, err := c.GetMetricByKey(ctx, exampleNotFoundMetricKey)
	if err != nil {
		fmt.Printf("Ошибка получения метрики %s: %v\n", notFoundMetric, err)
		return
	}
	fmt.Printf("Не найдена метрика Key: %s получено пустое значение %+v\n", exampleNotFoundMetricKey, notFoundMetric)

	// Output:
	// Найдена метрика Key: counter:name-1 Val: 123
	// Не найдена метрика Key: gauge:not-found-metric получено пустое значение <nil>
}
