package counter

import (
	"context"
	"github.com/aridae/go-metrics-store/internal/metrics-store-server/models"
	"time"
)

type metricsRepo interface {
	GetLatestState(ctx context.Context, metricType models.ScalarMetricType, metricName string) (*models.ScalarMetric, error)
	Save(ctx context.Context, metric models.ScalarMetric) error
}

type Handler struct {
	metricsRepo metricsRepo
	now         func() time.Time
}

func NewHandler(repo metricsRepo) *Handler {
	return &Handler{
		metricsRepo: repo,
		now:         func() time.Time { return time.Now().UTC() },
	}
}
