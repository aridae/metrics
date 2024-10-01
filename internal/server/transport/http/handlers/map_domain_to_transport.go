package handlers

import (
	"fmt"
	"github.com/aridae/go-metrics-store/internal/server/models"
	httpmodels "github.com/aridae/go-metrics-store/internal/server/transport/http/models"
	"strconv"
)

func buildMetricTransportModel(metric models.ScalarMetric) (httpmodels.Metric, error) {
	switch metric.Type() {
	case models.ScalarMetricTypeCounter:
		return buildCounterTransportModel(metric.Name(), metric.Value().String())
	case models.ScalarMetricTypeGauge:
		return buildGaugeTransportModel(metric.Name(), metric.Value().String())
	default:
		return httpmodels.Metric{}, fmt.Errorf("unknown metric type: %s", metric.Type())
	}
}

func buildCounterTransportModel(name, val string) (httpmodels.Metric, error) {
	typedVal, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return httpmodels.Metric{}, fmt.Errorf("provided value '%s' is not a valid integer", val)
	}

	return httpmodels.Metric{
		ID:    name,
		MType: counter,
		Delta: &typedVal,
	}, nil
}

func buildGaugeTransportModel(name, val string) (httpmodels.Metric, error) {
	typedVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return httpmodels.Metric{}, fmt.Errorf("provided value '%s' is not a valid float64", val)
	}

	return httpmodels.Metric{
		ID:    name,
		MType: gauge,
		Value: &typedVal,
	}, nil
}
