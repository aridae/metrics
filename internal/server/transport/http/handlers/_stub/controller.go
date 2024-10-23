package _stub

import (
	"context"
	"github.com/aridae/go-metrics-store/internal/server/models"
)

type ControllerNoErrStub struct{}

func (stub *ControllerNoErrStub) UpsertMetric(ctx context.Context, metricToRegister models.MetricUpsert) (models.Metric, error) {
	return models.Metric{}, nil
}

func (stub *ControllerNoErrStub) UpsertMetricsBatch(ctx context.Context, metricToRegister []models.MetricUpsert) ([]models.Metric, error) {
	return nil, nil
}

func (stub *ControllerNoErrStub) GetMetricByKey(ctx context.Context, metricKey models.MetricKey) (*models.Metric, error) {
	return nil, nil
}

func (stub *ControllerNoErrStub) GetAllMetrics(ctx context.Context) ([]models.Metric, error) {
	return nil, nil
}
