package inmemory

import (
	"context"
	"fmt"

	"github.com/aridae/go-metrics-store/internal/server/models"
	tsstorage "github.com/aridae/go-metrics-store/pkg/timeseries-storage"
)

func (r *repo) GetLatestState(ctx context.Context, key models.MetricKey) (*models.ScalarMetric, error) {
	val := r.storage.GetLatest(ctx, tsstorage.Key(key.String()))
	if val == nil {
		return nil, nil
	}

	scalar, ok := val.(models.ScalarMetric)
	if !ok {
		return nil, fmt.Errorf("unexpected type %T, expecting models.ScalarMetric type", val)
	}

	return &scalar, nil
}

func (r *repo) GetAllLatestStates(ctx context.Context) ([]models.ScalarMetric, error) {
	rawMetrics := r.storage.GetAllLatest(ctx)

	metrics := make([]models.ScalarMetric, 0, len(rawMetrics))

	for _, rawMetric := range rawMetrics {
		scalar, ok := rawMetric.(models.ScalarMetric)
		if !ok {
			return nil, fmt.Errorf("unexpected type %T, expecting models.ScalarMetric type", rawMetric)
		}

		metrics = append(metrics, scalar)
	}

	return metrics, nil
}
