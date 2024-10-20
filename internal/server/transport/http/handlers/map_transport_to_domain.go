package handlers

import (
	"errors"
	"fmt"
	"github.com/aridae/go-metrics-store/internal/server/models"
	httpmodels "github.com/aridae/go-metrics-store/internal/server/transport/http/models"
)

func buildMetricDomainModel(httpMetric httpmodels.Metric) (models.MetricUpsert, error) {
	factory, err := resolveMetricFactoryForMetricType(httpMetric.MType)
	if err != nil {
		return models.MetricUpsert{}, fmt.Errorf("failed to resolve metric factory for <type:%s> <id:%s>: %w", httpMetric.MType, httpMetric.ID, err)
	}

	transportMetricValue, err := resolveDomainMetricValueFromTransportModel(httpMetric)
	if err != nil {
		return models.MetricUpsert{}, fmt.Errorf("failed to resolve metric value from transport model: %w", err)
	}

	metricName := httpMetric.ID
	metricValue, err := factory.ParseMetricValue(transportMetricValue)
	if err != nil {
		return models.MetricUpsert{}, fmt.Errorf("failed to parse transport metric value: %w", err)
	}

	metric := factory.CreateMetricUpsert(metricName, metricValue)

	return metric, nil
}

func resolveDomainMetricValueFromTransportModel(metric httpmodels.Metric) (string, error) {
	switch metric.MType {
	case gauge:
		if metric.Value == nil {
			return "", errors.New("'value' field is required for gauge metric")
		}
		return *metric.Value, nil
	case counter:
		if metric.Delta == nil {
			return "", errors.New("'delta' field is required for counter metric")
		}
		return *metric.Delta, nil
	default:
		return "", fmt.Errorf("unknown metric type: %s", metric.MType)
	}
}
