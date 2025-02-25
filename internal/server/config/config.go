package config

import (
	"sync"
	"time"

	"github.com/aridae/go-metrics-store/pkg/logger"
)

const (
	addressDefaultVal          = "localhost:8080"
	storeIntervalDefault       = time.Duration(300) * time.Second
	fileStoragePathDefault     = "./.data"
	restoreDefault             = true
	databaseMaxOpenConnDefault = 5
)

var (
	once         sync.Once
	globalConfig *Config
)

type Config struct {
	CryptoKey           string
	Address             string
	FileStoragePath     string
	DatabaseDsn         string
	Key                 string
	StoreInterval       time.Duration
	DatabaseMaxOpenConn int
	Restore             bool
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

	flagsValues := parseFlags()

	envValues, err := parseEnv()
	if err != nil {
		logger.Errorf("error parsing environment, proceeding without env overrides: %v", err)
	}

	configFilePath := flagsValues.ConfigFilePath
	if configFilePath == "" && envValues.ConfigFilePath != nil {
		configFilePath = *envValues.ConfigFilePath
	}

	var jsonFileValues *jsonconf
	if configFilePath != "" {
		jsonFileValues, err = parseJSONFile(configFilePath)
		if err != nil {
			logger.Errorf("error parsing json config, proceeding without yaml overrides: %v", err)
		}
	}

	// json файл если есть, используем его - с наименьшим приоритетом
	if jsonFileValues != nil {
		jsonFileValues.override(c)
	}

	// перезатираем значениями, переданными через флаги
	flagsValues.override(c)

	// env, если есть, затирает флаги
	if envValues != nil {
		envValues.override(c)
	}
}

func (c *Config) defaults() {
	c.Address = addressDefaultVal
	c.StoreInterval = storeIntervalDefault
	c.FileStoragePath = fileStoragePathDefault
	c.Restore = restoreDefault
	c.DatabaseMaxOpenConn = databaseMaxOpenConnDefault
}

func (c *Config) overrideAddressIfNotDefault(address string, source string) {
	if address == addressDefaultVal {
		logger.Debugf("source %s provided default Address value, not overriding", source)
		return
	}

	logger.Infof("overriding Address from %s: (%s)-->(%s)", source, c.Address, address)
	c.Address = address
}

func (c *Config) overrideStoreIntervalIfNotDefault(storeInterval time.Duration, source string) {
	if storeInterval == storeIntervalDefault {
		logger.Debugf("source %s provided default StoreInterval value, not overriding", source)
		return
	}

	logger.Infof("overriding StoreInterval from %s: (%s)-->(%s)", source, c.StoreInterval, storeInterval)
	c.StoreInterval = storeInterval
}

func (c *Config) overrideFileStoragePathIfNotDefault(fileStoragePath string, source string) {
	if fileStoragePath == fileStoragePathDefault {
		logger.Debugf("source %s provided default FileStoragePath value, not overriding", source)
		return
	}

	logger.Infof("overriding FileStoragePath from %s: (%s)-->(%s)", source, c.FileStoragePath, fileStoragePath)
	c.FileStoragePath = fileStoragePath
}

func (c *Config) overrideRestoreIfNotDefault(restore bool, source string) {
	if restore {
		logger.Debugf("source %s provided default Restore value, not overriding", source)
		return
	}

	logger.Infof("overriding Restore from %s: (%t)-->(%t)", source, c.Restore, restore)
	c.Restore = restore
}

func (c *Config) overrideDatabaseDNSIfNotDefault(dns string, source string) {
	if dns == "" {
		logger.Debugf("source %s provided empty dns value, not overriding", source)
		return
	}

	logger.Infof("overriding dns from %s", source)
	c.DatabaseDsn = dns
}

func (c *Config) overrideKeyIfNotDefault(key string, source string) {
	if key == "" {
		logger.Debugf("source %s provided empty key value, not overriding", source)
		return
	}

	logger.Infof("overriding key from %s", source)
	c.Key = key
}

func (c *Config) overrideCryptoKeyIfNotDefault(cryptoKey string, source string) {
	if cryptoKey == "" {
		logger.Debugf("source %s provided empty crypto key value, not overriding", source)
		return
	}

	logger.Infof("overriding cryptoKey from %s", source)
	c.CryptoKey = cryptoKey
}
