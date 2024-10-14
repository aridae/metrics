package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
)

type environs struct {
	AddressOverride              *string `env:"ADDRESS"`
	StoreIntervalSecondsOverride *int64  `env:"STORE_INTERVAL"`
	FileStoragePathOverride      *string `env:"FILE_STORAGE_PATH"`
	RestoreOverride              *bool   `env:"RESTORE"`
}

func readEnv() (environs, error) {
	envs := environs{}

	err := env.Parse(&envs)
	if err != nil {
		return environs{}, fmt.Errorf("failed to parse env variables: %w", err)
	}

	return envs, nil
}

func (e environs) override(cfg *Config) {
	if e.AddressOverride != nil {
		cfg.overrideAddressIfNotDefault(*e.AddressOverride, "env")
	}

	if e.StoreIntervalSecondsOverride != nil {
		storeInterval := time.Duration(*e.StoreIntervalSecondsOverride) * time.Second
		cfg.overrideStoreIntervalIfNotDefault(storeInterval, "env")
	}

	if e.FileStoragePathOverride != nil {
		cfg.overrideFileStoragePathIfNotDefault(*e.FileStoragePathOverride, "env")
	}

	if e.RestoreOverride != nil {
		cfg.overrideRestoreIfNotDefault(*e.RestoreOverride, "env")
	}
}
