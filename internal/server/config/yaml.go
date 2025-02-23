package config

import (
	"github.com/aridae/go-metrics-store/pkg/logger"
	"gopkg.in/yaml.v3"
	"os"
)

func mustParseYaml(cnf *Config, path string) {
	file, err := os.Open(path)
	if err != nil {
		logger.Fatalf("failed to open config yaml file: %v", err)
	}
	defer file.Close()

	if err = yaml.NewDecoder(file).Decode(&cnf); err != nil {
		logger.Fatalf("failed to parse yaml file: %v", err)
	}
}
