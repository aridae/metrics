package config

import (
	"github.com/aridae/go-metrics-store/pkg/logger"
	goflags "github.com/jessevdk/go-flags"
)

func mustParseFlags(cnf *Config) {
	_, err := goflags.Parse(cnf)
	if err != nil {
		logger.Fatalf("error parsing command line flags: %v", err)
	}
}
