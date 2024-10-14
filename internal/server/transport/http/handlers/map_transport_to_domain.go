package handlers

import (
	"errors"
	"fmt"

	metricsfabrics "github.com/aridae/go-metrics-store/internal/server/metrics-fabrics"
	"github.com/aridae/go-metrics-store/internal/server/models"
	httpmodels "github.com/aridae/go-metrics-store/internal/server/transport/http/models"
)

func buildMetricDomainModel(
	transportMetric httpmodels.Metric,
	factory metricsfabrics.ScalarMetricFactory,
) (models.ScalarMetricToRegister, error) {
	transportMetricValue, err := resolveDomainMetricValueFromTransportModel(transportMetric)
	if err != nil {
		return models.ScalarMetricToRegister{}, fmt.Errorf("failed to resolve metric value from transport model: %w", err)
	}

	metricName := transportMetric.ID
	metricValue, err := factory.CastScalarMetricValue(transportMetricValue)
	if err != nil {
		return models.ScalarMetricToRegister{}, fmt.Errorf("failed to cast metric value: %w", err)
	}

	metric := factory.CreateScalarMetricToRegister(metricName, metricValue)

	return metric, nil
}

func resolveDomainMetricValueFromTransportModel(metric httpmodels.Metric) (any, error) {
	switch metric.MType {
	case gauge:
		if metric.Value == nil {
			return nil, errors.New("'value' field is required for gauge metric")
		}
		return *metric.Value, nil
	case counter:
		if metric.Delta == nil {
			return nil, errors.New("'delta' field is required for counter metric")
		}
		return *metric.Delta, nil
	default:
		return nil, fmt.Errorf("unknown metric type: %s", metric.MType)
	}
}
