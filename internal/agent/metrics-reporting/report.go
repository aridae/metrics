package metricsreporting

import (
	"context"
	"fmt"
	metricsservice "github.com/aridae/go-metrics-store/internal/agent/downstreams/metrics-service"
)

const (
	counterType = "counter"
	gaugeType   = "gauge"
)

func reportMetricWithJSONPayload(ctx context.Context, metricsService metricsService, pack metricsPack) error {
	metricsBatch := make([]metricsservice.Metric, 0, len(pack.gauges)+len(pack.counters))

	for gaugeName, gaugeVal := range pack.gauges {
		jsonGauge, err := buildMetricJSONPayload(gaugeType, gaugeName, gaugeVal)
		if err != nil {
			return fmt.Errorf("failed to build gauge <name:%s> json-serializable struct: %w", gaugeName, err)
		}

		metricsBatch = append(metricsBatch, jsonGauge)
	}

	for counterName, counterVal := range pack.counters {
		jsonCounter, err := buildMetricJSONPayload(counterType, counterName, counterVal)
		if err != nil {
			return fmt.Errorf("failed to build counter <name:%s> json-serializable struct: %w", counterName, err)
		}

		metricsBatch = append(metricsBatch, jsonCounter)
	}

	if err := metricsService.UpdateMetricsBatch(ctx, metricsBatch); err != nil {
		return fmt.Errorf("failed to update metrics batch: %w", err)
	}

	return nil
}

func buildMetricJSONPayload(
	mtype string,
	name string,
	val any,
) (metricsservice.Metric, error) {
	switch mtype {
	case counterType:
		counterVal, ok := val.(counter)
		if !ok {
			return metricsservice.Metric{}, fmt.Errorf("value is not int64")
		}
		int64Val := int64(counterVal)
		return metricsservice.Metric{
			ID:    name,
			MType: mtype,
			Delta: &int64Val,
		}, nil
	case gaugeType:
		gaugeVal, ok := val.(gauge)
		if !ok {
			return metricsservice.Metric{}, fmt.Errorf("value is not float64")
		}
		float64Val := float64(gaugeVal)
		return metricsservice.Metric{
			ID:    name,
			MType: mtype,
			Value: &float64Val,
		}, nil
	default:
		return metricsservice.Metric{}, fmt.Errorf("unsupported metric type: %s", mtype)
	}
}
