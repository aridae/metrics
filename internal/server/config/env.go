package config

import (
	"github.com/caarlos0/env/v6"
	"log"
	"time"
)

type environs struct {
	AddressOverride              string `env:"ADDRESS"`
	StoreIntervalSecondsOverride int64  `env:"STORE_INTERVAL"`
	FileStoragePathOverride      string `env:"FILE_STORAGE_PATH"`
	RestoreOverride              bool   `env:"RESTORE"`
}

func readEnv() environs {
	envs := environs{}

	err := env.Parse(&envs)
	if err != nil {
		log.Fatalf("failed to parse env variables: %v", err)
	}

	return envs
}

func (e environs) override(cfg *Config) {
	if e.AddressOverride != "" {
		log.Printf("overriding Address with env: %s", e.AddressOverride)
		cfg.Address = e.AddressOverride
	}

	if e.StoreIntervalSecondsOverride != 0 {
		log.Printf("overriding StoreIntervalSecondsOverride with env: %d", e.StoreIntervalSecondsOverride)
		cfg.StoreInterval = time.Duration(e.StoreIntervalSecondsOverride) * time.Second
	}

	if e.FileStoragePathOverride != "" {
		log.Printf("overriding FileStoragePath with env: %s", e.FileStoragePathOverride)
		cfg.FileStoragePath = e.FileStoragePathOverride
	}

	if e.RestoreOverride {
		log.Printf("overriding Restore with env: %t", e.RestoreOverride)
		cfg.Restore = e.RestoreOverride
	}
}
