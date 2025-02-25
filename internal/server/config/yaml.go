package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type yamls struct {
	AddressOverride         *string        `yaml:"address"`
	StoreIntervalOverride   *time.Duration `yaml:"store_interval"`
	FileStoragePathOverride *string        `yaml:"file_storage_path"`
	RestoreOverride         *bool          `yaml:"restore"`
}

func parseYaml(path string) (yamls, error) {
	file, err := os.Open(path)
	if err != nil {
		return yamls{}, fmt.Errorf("failed to open config yaml file: %w", err)
	}
	defer file.Close()

	yamlCnf := yamls{}
	if err = yaml.NewDecoder(file).Decode(&yamlCnf); err != nil {
		return yamls{}, fmt.Errorf("failed to parse yaml file: %w", err)
	}

	return yamlCnf, nil
}

func (f yamls) override(cfg *Config) {
	if f.AddressOverride != nil {
		cfg.overrideAddressIfNotDefault(*f.AddressOverride, "yamls config")
	}

	if f.StoreIntervalOverride != nil {
		cfg.overrideStoreIntervalIfNotDefault(*f.StoreIntervalOverride, "yamls config")
	}

	if f.FileStoragePathOverride != nil {
		cfg.overrideFileStoragePathIfNotDefault(*f.FileStoragePathOverride, "yamls config")
	}

	if f.RestoreOverride != nil {
		cfg.overrideRestoreIfNotDefault(*f.RestoreOverride, "yamls config")
	}
}
