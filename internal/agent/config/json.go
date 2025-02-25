package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type jsonconf struct {
	AddressOverride        *string        `json:"address"`
	ReportIntervalOverride *time.Duration `json:"report_interval"`
	PollIntervalOverride   *time.Duration `json:"poll_interval"`
	CryptoKey              *string        `env:"crypto_key"`
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

	if f.ReportIntervalOverride != nil {
		cfg.overrideReportIntervalIfNotDefault(*f.ReportIntervalOverride, "json config")
	}

	if f.PollIntervalOverride != nil {
		cfg.overridePollIntervalIfNotDefault(*f.PollIntervalOverride, "json config")
	}

	if f.CryptoKey != nil {
		cfg.overrideCryptoKeyIfNotDefault(*f.CryptoKey, "json config")
	}
}
