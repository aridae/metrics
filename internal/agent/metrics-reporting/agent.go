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
	metricsService metricsService

	reportersPoolSize int64

	pollInterval   time.Duration
	reportInterval time.Duration
}

func NewAgent(
	metricsService metricsService,
	pollInterval time.Duration,
	reportInterval time.Duration,
	reportersPoolSize int64,
) *Agent {
	return &Agent{
		metricsService:    metricsService,
		pollInterval:      pollInterval,
		reportInterval:    reportInterval,
		reportersPoolSize: reportersPoolSize,
	}
}

func (a *Agent) Run(ctx context.Context) {
	metricsQueue := make(chan metricsPack, 1000)
	go func() {
		<-ctx.Done()
		logger.Errorf("stopping agent due to context cancellation: %v", ctx.Err())
		close(metricsQueue)
	}()

	// (2) При этом количество одновременно исходящих запросов на сервер нужно ограничивать «сверху»
	for i := range a.reportersPoolSize {
		go a.runReportingWorkerLoop(ctx, i, metricsQueue)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	// (1) Перепланируйте архитектуру агента таким образом, чтобы сбор метрик (опрос runtime)
	// и их отправка осуществлялись в разных горутинах
	go func() {
		defer wg.Done()
		a.runPollingRuntimeLoop(ctx, metricsQueue)
	}()

	// (3) Добавьте ещё одну горутину, которая будет использовать пакет gopsutil и собирать дополнительные метрики
	go func() {
		defer wg.Done()
		a.runPollingGopsutilLoop(ctx, metricsQueue)
	}()

	wg.Wait()
}

func (a *Agent) runPollingRuntimeLoop(ctx context.Context, metricsQueue chan<- metricsPack) {
	pollTick := time.NewTicker(a.pollInterval)
	for {
		select {
		case <-ctx.Done():
			logger.Infof("stopping runPollingRuntimeLoop routine due to context cancellation")
			return
		case <-pollTick.C:
			logger.Infof("starting PollingRuntime routine <now:%s>", time.Now().UTC())

			pack := pollRuntime(ctx)

			queueMetricsPack(ctx, pack, metricsQueue)
		}
	}
}

func (a *Agent) runPollingGopsutilLoop(ctx context.Context, metricsQueue chan<- metricsPack) {
	pollTick := time.NewTicker(a.pollInterval)
	for {
		select {
		case <-ctx.Done():
			logger.Infof("stopping runPollingGopsutilLoop routine due to context cancellation")
			return
		case <-pollTick.C:
			logger.Infof("starting PollingGopsutil routine <now:%s>", time.Now().UTC())

			pack, err := pollGopsutil(ctx)
			if err != nil {
				logger.Errorf("error while polling gopsutil: %v", err)
				continue
			}

			queueMetricsPack(ctx, pack, metricsQueue)
		}
	}
}

func (a *Agent) runReportingWorkerLoop(ctx context.Context, reporterID int64, metricsQueue <-chan metricsPack) {
	reportTick := time.NewTicker(a.reportInterval)
	for {
		select {
		case <-ctx.Done():
			logger.Infof("stopping reporter worker #%d due to context cancellation", reporterID)
			return
		case <-reportTick.C:
			pack, err := dequeueMetricsPack(ctx, metricsQueue)
			if err != nil {
				logger.Errorf("reporter worker #%d failed dequeue metrics to report: %v", reporterID, err)
				continue
			}

			logger.Infof("reporter worker #%d received metrics pack to report", reporterID)
			if err = reportMetricWithJSONPayload(ctx, a.metricsService, pack); err != nil {
				logger.Errorf("reporter worker #%d failed to report metrics: %v", reporterID, err)
			}
		}
	}
}

func queueMetricsPack(ctx context.Context, pack metricsPack, metricsQueue chan<- metricsPack) {
	select {
	case metricsQueue <- pack:
	case <-ctx.Done():
		logger.Infof("terminating queueMetricsPack due to context cancellation")
		return
	}
}

func dequeueMetricsPack(ctx context.Context, metricsQueue <-chan metricsPack) (metricsPack, error) {
	select {
	case pack := <-metricsQueue:
		return pack, nil
	case <-ctx.Done():
		logger.Infof("terminating dequeueMetricsPack due to context cancellation")
		return metricsPack{}, ctx.Err()
	}
}
