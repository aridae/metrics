package _stub

import (
	"context"
	metricsupsertstrategies "github.com/aridae/go-metrics-store/internal/server/metrics-upsert-strategies"
	"github.com/aridae/go-metrics-store/internal/server/models"
)

type ControllerNoErrStub struct{}

func (stub *ControllerNoErrStub) UpsertScalarMetric(ctx context.Context, metricToRegister models.ScalarMetricToRegister, strategy metricsupsertstrategies.Strategy) (models.ScalarMetric, error) {
	return models.ScalarMetric{}, nil
}

func (stub *ControllerNoErrStub) GetScalarMetricLatestState(ctx context.Context, metricKey models.MetricKey) (*models.ScalarMetric, error) {
	return nil, nil
}

func (stub *ControllerNoErrStub) GetAllScalarMetricsLatestStates(ctx context.Context) ([]models.ScalarMetric, error) {
	return nil, nil
}
