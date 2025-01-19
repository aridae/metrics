package metricinmemrepo

import (
	"context"

	"github.com/aridae/go-metrics-store/internal/server/models"
	"github.com/aridae/go-metrics-store/pkg/inmem"
)

func (r *repo) Save(ctx context.Context, metric models.Metric) error {
	key := metric.GetKey().String()

	r.db.Save(ctx, inmem.Key(key), metric)

	return nil
}
