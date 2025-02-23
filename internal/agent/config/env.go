package config

import (
	"github.com/aridae/go-metrics-store/pkg/logger"
	"github.com/caarlos0/env/v6"
)

func mustParseEnv(cnf *Config) {
	err := env.Parse(cnf)
	if err != nil {
		logger.Fatalf("error parsing environment: %v", err)
	}
}
