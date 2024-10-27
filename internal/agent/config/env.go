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
}
