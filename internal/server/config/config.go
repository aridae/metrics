package config

import (
	"sync"
)

const (
	yamlConfigPath = "./config/server.yaml"
)

var (
	once         sync.Once
	globalConfig *Config
)

type Config struct {
	Address              string `short:"a" description:"Address of server" env:"ADDRESS" yaml:"address"`
	FileStoragePath      string `short:"f" description:"Backup file path" env:"FILE_STORAGE_PATH" yaml:"file_storage_path"`
	DatabaseDsn          string `short:"d" description:"Database DSN" env:"DATABASE_DSN"`
	Key                  string `short:"k" description:"Ключ для подписания запросов SHA256 подписью" env:"KEY"`
	StoreIntervalSeconds int    `short:"i" description:"Backup store interval" env:"STORE_INTERVAL" yaml:"store_interval"`
	Restore              bool   `short:"r" description:"Restore from backup file on start" env:"RESTORE" yaml:"restore"`
	CryptoKey            string `long:"crypto-key" description:"путь до файла с приватным ключом" env:"CRYPTO_KEY"`
	DatabaseMaxOpenConn  int
}

func Obtain() *Config {
	once.Do(func() {
		globalConfig = &Config{}
		globalConfig.init()
	})

	return globalConfig
}

var defaultConfig = Config{
	Address:              "localhost:8080",
	FileStoragePath:      "./.data",
	StoreIntervalSeconds: 300,
	Restore:              true,
}

func (cnf *Config) init() {
	*cnf = defaultConfig

	// инициализация структуры конфига из yaml файла
	mustParseYaml(cnf, yamlConfigPath)

	// перезатираем значениями, переданными через флаги
	mustParseFlags(cnf)

	// env, если есть, затирает флаги
	mustParseEnv(cnf)
}
