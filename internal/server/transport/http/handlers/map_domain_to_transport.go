package handlers

import (
	"fmt"

	"github.com/aridae/go-metrics-store/internal/server/models"
	httpmodels "github.com/aridae/go-metrics-store/internal/server/transport/http/models"
	"github.com/aridae/go-metrics-store/pkg/pointer"
)

func buildMetricTransportModel(metric models.Metric) (httpmodels.Metric, error) {
	switch metric.GetType() {
	case models.MetricTypeCounter:
		return buildCounterTransportModel(metric.GetName(), metric.GetValue()), nil
	case models.MetricTypeGauge:
		return buildGaugeTransportModel(metric.GetName(), metric.GetValue()), nil
	default:
		return httpmodels.Metric{}, fmt.Errorf("unknown metric type: %s", metric.GetType())
	}
}

func buildCounterTransportModel(name string, val models.MetricValue) httpmodels.Metric {
	return httpmodels.Metric{
		ID:    name,
		MType: counter,
		Delta: pointer.To(val.UnsafeCastInt()),
	}
}

func buildGaugeTransportModel(name string, val models.MetricValue) httpmodels.Metric {
	return httpmodels.Metric{
		ID:    name,
		MType: gauge,
		Value: pointer.To(val.UnsafeCastFloat()),
	}
}
