package gauge

import (
	"context"
	"github.com/aridae/go-metrics-store/internal/metrics-store-server/models"
	"time"
)

type metricsRepo interface {
	Save(ctx context.Context, metric models.ScalarMetric) error
}

type Handler struct {
	repo metricsRepo
	now  func() time.Time
}

func NewHandler(repo metricsRepo) *Handler {
	return &Handler{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}
