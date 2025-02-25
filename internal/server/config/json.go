package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type jsonconf struct {
	AddressOverride         *string        `json:"address"`
	StoreIntervalOverride   *time.Duration `json:"store_interval"`
	FileStoragePathOverride *string        `json:"file_storage_path"`
	RestoreOverride         *bool          `json:"restore"`
}

func parseJSONFile(path string) (*jsonconf, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config json file: %w", err)
	}
	defer file.Close()

	cnf := jsonconf{}
	if err = json.NewDecoder(file).Decode(&cnf); err != nil {
		return nil, fmt.Errorf("failed to parse json file: %w", err)
	}

	return &cnf, nil
}

func (f jsonconf) override(cfg *Config) {
	if f.AddressOverride != nil {
		cfg.overrideAddressIfNotDefault(*f.AddressOverride, "json config")
	}

	if f.StoreIntervalOverride != nil {
		cfg.overrideStoreIntervalIfNotDefault(*f.StoreIntervalOverride, "json config")
	}

	if f.FileStoragePathOverride != nil {
		cfg.overrideFileStoragePathIfNotDefault(*f.FileStoragePathOverride, "json config")
	}

	if f.RestoreOverride != nil {
		cfg.overrideRestoreIfNotDefault(*f.RestoreOverride, "json config")
	}
}
