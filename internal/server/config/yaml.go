package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

type yamls struct {
	AddressOverride         string        `yaml:"address"`
	StoreIntervalOverride   time.Duration `yaml:"store_interval"`
	FileStoragePathOverride string        `yaml:"file_storage_path"`
	RestoreOverride         bool          `yaml:"restore"`
}

func parseYaml(path string) yamls {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open config yaml file: %v", err)
	}

	yamlCnf := yamls{}

	if err = yaml.NewDecoder(file).Decode(&yamlCnf); err != nil {
		log.Fatalf("failed to parse yaml file: %v", err)
	}

	return yamlCnf
}

func (f yamls) override(cfg *Config) {
	if f.AddressOverride != "" {
		log.Printf("overriding Address with yaml values: %s", f.AddressOverride)
		cfg.Address = f.AddressOverride
	}

	if f.StoreIntervalOverride != -1 {
		log.Printf("overriding StoreInterval with yaml values: %s", f.StoreIntervalOverride)
		cfg.StoreInterval = f.StoreIntervalOverride
	}

	if f.FileStoragePathOverride != "" {
		log.Printf("overriding FileStoragePath with yaml values: %s", f.FileStoragePathOverride)
		cfg.FileStoragePath = f.FileStoragePathOverride
	}

	if f.RestoreOverride {
		log.Printf("overriding Restore with yaml values: %t", f.RestoreOverride)
		cfg.Restore = f.RestoreOverride
	}
}
