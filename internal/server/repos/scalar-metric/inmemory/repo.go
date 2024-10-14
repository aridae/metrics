package inmemory

import (
	"context"
	"fmt"

	"github.com/aridae/go-metrics-store/internal/server/logger"
	scalarmetric "github.com/aridae/go-metrics-store/internal/server/repos/scalar-metric"
	tsstorage "github.com/aridae/go-metrics-store/pkg/timeseries-storage"
)

type timeseriesStorage interface {
	Save(ctx context.Context, key tsstorage.Key, value tsstorage.TimeseriesValue)
	GetLatest(ctx context.Context, key tsstorage.Key) tsstorage.TimeseriesValue
	GetAllLatest(ctx context.Context) []tsstorage.TimeseriesValue
}

type repo struct {
	storage timeseriesStorage
}

func NewRepositoryImplementation(storage timeseriesStorage) scalarmetric.Repository {
	return &repo{storage: storage}
}

func (r *repo) Healthcheck(_ context.Context) error {
	if r == nil {
		return fmt.Errorf("nil repo receiver")
	}

	logger.Obtain().Infof("inmemory.repo is rather healthy!")
	return nil
}
