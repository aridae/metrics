package metricinmemrepo

import (
	"context"
	"fmt"

	"github.com/aridae/go-metrics-store/internal/server/models"
	"github.com/aridae/go-metrics-store/pkg/inmem"
)

func (r *repo) GetByKey(ctx context.Context, key models.MetricKey) (*models.Metric, error) {
	val := r.db.GetLatest(ctx, inmem.Key(key.String()))
	if val == nil {
		return nil, nil
	}

	metric, ok := val.(models.Metric)
	if !ok {
		return nil, fmt.Errorf("unexpected type %T, expecting models.Metric type", val)
	}

	return &metric, nil
}

func (r *repo) GetAll(ctx context.Context) ([]models.Metric, error) {
	rawMetrics := r.db.GetAllLatest(ctx)

	metrics := make([]models.Metric, 0, len(rawMetrics))

	for _, rawMetric := range rawMetrics {
		metric, ok := rawMetric.(models.Metric)
		if !ok {
			return nil, fmt.Errorf("unexpected type %T, expecting models.Metric type", rawMetric)
		}

		metrics = append(metrics, metric)
	}

	return metrics, nil
}
