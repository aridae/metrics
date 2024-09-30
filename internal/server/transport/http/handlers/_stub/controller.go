package _stub

import (
	"context"
	"github.com/aridae/go-metrics-store/internal/server/models"
)

type ControllerNoErrStub struct{}

func (stub *ControllerNoErrStub) UpsertScalarMetric(_ context.Context, _ models.ScalarMetricToRegister) error {
	return nil
}

func (stub *ControllerNoErrStub) GetScalarMetricLatestState(ctx context.Context, metricKey models.MetricKey) (*models.ScalarMetric, error) {
	return nil, nil
}

func (stub *ControllerNoErrStub) GetAllScalarMetricsLatestStates(ctx context.Context) ([]models.ScalarMetric, error) {
	return nil, nil
}
