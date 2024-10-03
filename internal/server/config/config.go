package config

import (
	"sync"
	"time"
)

const (
	yamlConfigPath = "./config/server.yaml"

	addressDefaultVal      = "localhost:8080"
	storeIntervalDefault   = time.Duration(300) * time.Second
	fileStoragePathDefault = "./.data"
	restoreDefault         = true
)

var (
	once         sync.Once
	globalConfig *Config
)

type Config struct {
	Address         string
	StoreInterval   time.Duration
	FileStoragePath string
	Restore         bool
}

func Obtain() *Config {
	once.Do(func() {
		globalConfig = &Config{}
		globalConfig.init()
	})

	return globalConfig
}

func (c *Config) init() {
	c.defaults()

	// инициализация структуры конфига из yaml файла
	parseYaml(yamlConfigPath).override(c)

	// перезатираем значениями, переданными через флаги
	parseFlags().override(c)

	// env, если есть, затирает флаги
	readEnv().override(c)
}

func (c *Config) defaults() {
	c.Address = addressDefaultVal
	c.StoreInterval = storeIntervalDefault
	c.FileStoragePath = fileStoragePathDefault
	c.Restore = restoreDefault
}
