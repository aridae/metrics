package metric

import (
	"context"
	"github.com/aridae/go-metrics-store/internal/metrics-store-server/models"
	tsstorage "github.com/aridae/go-metrics-store/pkg/ts-storage"
)

func (r *Repository) Save(ctx context.Context, metric models.ScalarMetric) error {
	key := string(metric.Type) + ":" + metric.Name

	r.storage.Save(ctx, tsstorage.Key(key), metric)

	return nil
}
