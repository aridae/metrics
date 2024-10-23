package handlers

import (
	"fmt"
	metricsfabrics "github.com/aridae/go-metrics-store/internal/server/models/metrics-fabrics"
)

func resolveMetricFactoryForMetricType(metricType string) (metricsfabrics.MetricFactory, error) {
	switch metricType {
	case counter:
		return metricsfabrics.ObtainCounterMetricFactory(), nil
	case gauge:
		return metricsfabrics.ObtainGaugeMetricFactory(), nil
	default:
		return nil, fmt.Errorf("unknown metric type: %s", metricType)
	}
}
