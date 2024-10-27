package metricsreporting

import (
	"context"
	metricsservice "github.com/aridae/go-metrics-store/internal/agent/downstreams/metrics-service"
	"github.com/aridae/go-metrics-store/pkg/logger"
	"sync"
	"time"
)

type metricsService interface {
	UpdateMetricsBatch(_ context.Context, metrics []metricsservice.Metric) error
}

type Agent struct {
	mu sync.RWMutex

	gauges   map[string]gauge
	counters map[string]counter

	metricsService metricsService

	pollInterval   time.Duration
	reportInterval time.Duration
}

func NewAgent(
	metricsService metricsService,
	pollInterval time.Duration,
	reportInterval time.Duration,
) *Agent {
	return &Agent{
		gauges:         make(map[string]gauge),
		counters:       make(map[string]counter),
		metricsService: metricsService,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
	}
}

func (a *Agent) Run(ctx context.Context) {
	pollTick := time.NewTicker(a.pollInterval)
	reportTick := time.NewTicker(a.reportInterval)

	for {
		select {
		case <-ctx.Done():
			logger.Obtain().Infof("stopping agent due to context cancellation")
			return
		case <-pollTick.C:
			logger.Obtain().Infof("starting metrics polling routine <now:%s>", time.Now().UTC())
			a.poll(ctx)
		case <-reportTick.C:
			logger.Obtain().Infof("starting metrics reporting routine <now:%s>", time.Now().UTC())
			if err := a.reportMetricWithJSONPayload(ctx); err != nil {
				logger.Obtain().Errorf("error reporting metrics: %v", err)
			}
		}
	}
}
