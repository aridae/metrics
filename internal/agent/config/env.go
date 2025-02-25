package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

func parseEnv(cnf *Config) {
	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		cnf.Address = envAddress
	}

	if envKey := os.Getenv("KEY"); envKey != "" {
		cnf.Key = envKey
	}

	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		cnf.CryptoKey = envCryptoKey
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		reportIntervalSec, err := strconv.ParseInt(envReportInterval, 10, 64)
		if err != nil {
			log.Fatalf("invalid REPORT_INTERVAL environment variable, int64 value expected: %v", err)
		}
		cnf.ReportInterval = time.Duration(reportIntervalSec) * time.Second
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		pollIntervalSec, err := strconv.ParseInt(envPollInterval, 10, 64)
		if err != nil {
			log.Fatalf("invalid POLL_INTERVAL environment variable, int64 value expected: %v", err)
		}
		cnf.PollInterval = time.Duration(pollIntervalSec) * time.Second
	}

	if envReportersPoolSize := os.Getenv("RATE_LIMIT"); envReportersPoolSize != "" {
		reportersPoolSize, err := strconv.ParseInt(envReportersPoolSize, 10, 64)
		if err != nil {
			log.Fatalf("invalid RATE_LIMIT environment variable, int64 value expected: %v", err)
		}
		cnf.ReportersPoolSize = reportersPoolSize
	}
}
