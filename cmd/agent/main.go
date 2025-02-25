package main

import (
	"context"
	"crypto/rsa"
	"github.com/aridae/go-metrics-store/internal/agent/config"
	metricsservice "github.com/aridae/go-metrics-store/internal/agent/downstreams/metrics-service"
	metricsreporting "github.com/aridae/go-metrics-store/internal/agent/metrics-reporting"
	rsamw "github.com/aridae/go-metrics-store/internal/server/transport/http/mw/rsa-mw"
	"github.com/aridae/go-metrics-store/internal/server/transport/http/mw/sha256-mw"
	"github.com/aridae/go-metrics-store/pkg/logger"
	rsacrypto "github.com/aridae/go-metrics-store/pkg/rsa-crypto"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	cnf := config.Obtain()

	logger.Infof("Starting Agent app with build flags:\n\nBuild version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	clientMiddlewares := []func(http.RoundTripper) http.RoundTripper{
		sha256mw.SignRequestClientMiddleware(cnf.Key),
	}

	if cnf.CryptoKey != "" {
		pubKey := mustParsePublicKey(cnf.CryptoKey)
		clientMiddlewares = append(clientMiddlewares, rsamw.EncryptRequestClientMiddleware(pubKey))
	}

	metricsServiceClient := metricsservice.NewClient(cnf.Address, clientMiddlewares...)

	metricsAgent := metricsreporting.NewAgent(
		metricsServiceClient,
		cnf.PollInterval,
		cnf.ReportInterval,
		cnf.ReportersPoolSize,
	)

	metricsAgent.Run(ctx)
}

func mustParsePublicKey(path string) *rsa.PublicKey {
	pubKey, err := rsacrypto.FromFile(path, rsacrypto.ParsePublicKey)
	if err != nil {
		logger.Fatalf("failed to load public key: %v", err)
	}

	return pubKey
}
