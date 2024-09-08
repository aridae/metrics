package metric

import (
	"context"
	"github.com/aridae/go-metrics-store/internal/server/models"
	tsstorage "github.com/aridae/go-metrics-store/pkg/ts-storage"
)

func (r *Repository) GetLatestState(ctx context.Context, metricType models.ScalarMetricType, metricName string) (*models.ScalarMetric, error) {
	key := string(metricType) + ":" + metricName

	val := r.storage.GetLatest(ctx, tsstorage.Key(key))

	if val == nil {
		return nil, nil
	}

	scalar, _ := val.(*models.ScalarMetric)

	return scalar, nil
}
