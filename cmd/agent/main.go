package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/aridae/go-metrics-store/internal/agent/config"
	metricsservice "github.com/aridae/go-metrics-store/internal/agent/downstreams/metrics-service"
	metricsreporting "github.com/aridae/go-metrics-store/internal/agent/metrics-reporting"
	"github.com/aridae/go-metrics-store/internal/server/transport/http/mw/sha256-mw"
	"github.com/aridae/go-metrics-store/pkg/logger"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)

		<-signalCh

		logger.Infof("Got signal, shutting down...")

		// If you fail to cancel the context, the goroutine that WithCancel or WithTimeout created
		// will be retained in memory indefinitely (until the program shuts down), causing a memory leak.
		cancel()
	}()

	logger.Infof("Starting Agent app with build flags:\n\nBuild version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	cnf := config.Init()

	metricsServiceClient := metricsservice.NewClient(cnf.Address,
		sha256mw.SignRequestClientMiddleware(cnf.Key),
	)
	metricsAgent := metricsreporting.NewAgent(metricsServiceClient, cnf.PollInterval, cnf.ReportInterval, cnf.ReportersPoolSize)

	metricsAgent.Run(ctx)
}
