package config

import (
	"github.com/aridae/go-metrics-store/pkg/logger"
	"sync"
	"time"
)

var (
	once         sync.Once
	globalConfig *Config
)

type Config struct {
	CryptoKey         string
	Address           string
	Key               string
	PollInterval      time.Duration
	ReportInterval    time.Duration
	ReportersPoolSize int64
}

var defaultConfig = Config{
	PollInterval:      2 * time.Second,
	ReportInterval:    10 * time.Second,
	Address:           "localhost:8080",
	Key:               "",
	ReportersPoolSize: 100,
}

func Obtain() *Config {
	once.Do(func() {
		globalConfig = &Config{}
		globalConfig.init()
	})

	return globalConfig
}

func (c *Config) init() {
	*c = defaultConfig

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
			logger.Errorf("error parsing json config, proceeding without json overrides: %v", err)
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

func (c *Config) overrideAddressIfNotDefault(address string, source string) {
	if address == defaultConfig.Address {
		logger.Debugf("source %s provided default Address value, not overriding", source)
		return
	}

	logger.Infof("overriding Address from %s: (%s)-->(%s)", source, c.Address, address)
	c.Address = address
}

func (c *Config) overridePollIntervalIfNotDefault(pollInterval time.Duration, source string) {
	if pollInterval == defaultConfig.PollInterval {
		logger.Debugf("source %s provided default PollInterval value, not overriding", source)
		return
	}

	logger.Infof("overriding PollInterval from %s: (%s)-->(%s)", source, c.PollInterval, pollInterval)
	c.PollInterval = pollInterval
}

func (c *Config) overrideReportIntervalIfNotDefault(reportInterval time.Duration, source string) {
	if reportInterval == defaultConfig.ReportInterval {
		logger.Debugf("source %s provided default ReportInterval value, not overriding", source)
		return
	}

	logger.Infof("overriding ReportInterval from %s: (%s)-->(%s)", source, c.ReportInterval, reportInterval)
	c.ReportInterval = reportInterval
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

func (c *Config) overridePoolSizeIfNotDefault(size int64, source string) {
	if size == defaultConfig.ReportersPoolSize {
		logger.Debugf("source %s provided default pool size value, not overriding", source)
		return
	}

	logger.Infof("overriding pool size from %s", source)
	c.ReportersPoolSize = size
}
