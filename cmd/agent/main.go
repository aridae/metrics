package main

import (
	"context"
	"github.com/aridae/go-metrics-store/internal/agent/config"
	metricsservice "github.com/aridae/go-metrics-store/internal/agent/downstreams/metrics-service"
	metricsreporting "github.com/aridae/go-metrics-store/internal/agent/metrics-reporting"
	"github.com/aridae/go-metrics-store/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)

		<-signalCh

		logger.Obtain().Info("Got signal, shutting down...")

		// If you fail to cancel the context, the goroutine that WithCancel or WithTimeout created
		// will be retained in memory indefinitely (until the program shuts down), causing a memory leak.
		cancel()
	}()

	cnf := config.Init()

	metricsServiceClient := metricsservice.NewClient(cnf.Address)
	metricsAgent := metricsreporting.NewAgent(metricsServiceClient)

	metricsAgent.Run(ctx)
}
