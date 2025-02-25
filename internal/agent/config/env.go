package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"time"
)

type envconf struct {
	ConfigFilePath        *string `env:"CONFIG"`
	CryptoKey             *string `env:"CRYPTO_KEY"`
	Address               *string `env:"ADDRESS"`
	Key                   *string `env:"KEY"`
	PollIntervalSeconds   *int64  `env:"POLL_INTERVAL"`
	ReportIntervalSeconds *int64  `env:"REPORT_INTERVAL"`
	ReportersPoolSize     *int64  `env:"RATE_LIMIT"`
}

func parseEnv() (*envconf, error) {
	envs := envconf{}

	err := env.Parse(&envs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse env variables: %w", err)
	}

	return &envs, nil
}

func (e envconf) override(cfg *Config) {
	if e.Address != nil {
		cfg.overrideAddressIfNotDefault(*e.Address, "env")
	}

	if e.ReportIntervalSeconds != nil {
		cfg.overrideReportIntervalIfNotDefault(time.Duration(*e.ReportIntervalSeconds)*time.Second, "env")
	}

	if e.PollIntervalSeconds != nil {
		cfg.overridePollIntervalIfNotDefault(time.Duration(*e.PollIntervalSeconds)*time.Second, "env")
	}

	if e.CryptoKey != nil {
		cfg.overrideCryptoKeyIfNotDefault(*e.CryptoKey, "env")
	}

	if e.Key != nil {
		cfg.overrideKeyIfNotDefault(*e.Key, "env")
	}

	if e.ReportersPoolSize != nil {
		cfg.overridePoolSizeIfNotDefault(*e.ReportersPoolSize, "env")
	}
}
