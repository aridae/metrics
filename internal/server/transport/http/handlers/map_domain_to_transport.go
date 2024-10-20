package handlers

import (
	"fmt"
	"github.com/aridae/go-metrics-store/internal/server/models"
	httpmodels "github.com/aridae/go-metrics-store/internal/server/transport/http/models"
)

func buildMetricTransportModel(metric models.Metric) (httpmodels.Metric, error) {
	switch metric.GetType() {
	case models.MetricTypeCounter:
		return buildCounterTransportModel(metric.GetName(), metric.GetValue().String()), nil
	case models.MetricTypeGauge:
		return buildGaugeTransportModel(metric.GetName(), metric.GetValue().String()), nil
	default:
		return httpmodels.Metric{}, fmt.Errorf("unknown metric type: %s", metric.GetType())
	}
}

func buildCounterTransportModel(name, val string) httpmodels.Metric {
	return httpmodels.Metric{
		ID:    name,
		MType: counter,
		Delta: &val,
	}
}

func buildGaugeTransportModel(name, val string) httpmodels.Metric {
	return httpmodels.Metric{
		ID:    name,
		MType: gauge,
		Value: &val,
	}
}
