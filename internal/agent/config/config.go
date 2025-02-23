package config

import "sync"

type Config struct {
	Address               string `short:"a" description:"адрес эндпоинта HTTP-сервера" env:"ADDRESS"`
	Key                   string `short:"k" description:"ключ для подписания запросов SHA256 подписью" env:"KEY"`
	PollIntervalSeconds   int64  `short:"p" description:"частота опроса метрик из пакета runtime" env:"POLL_INTERVAL"`
	ReportIntervalSeconds int64  `short:"r" description:"частота отправки метрик на сервер" env:"REPORT_INTERVAL"`
	ReportersPoolSize     int64  `short:"l" description:"количество одновременно исходящих запросов на сервер" env:"RATE_LIMIT"`
}

var (
	once         sync.Once
	globalConfig *Config
)

func Obtain() *Config {
	once.Do(func() {
		globalConfig = &Config{}
		globalConfig.init()
	})

	return globalConfig
}

var defaultConfig = Config{
	PollIntervalSeconds:   2,
	ReportIntervalSeconds: 10,
	Address:               "localhost:8080",
	ReportersPoolSize:     100,
}

func (cnf *Config) init() {
	*cnf = defaultConfig
	mustParseFlags(cnf)
	mustParseEnv(cnf)
}
