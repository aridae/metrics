package config

import (
	"flag"
	"time"
)

func parseFlags(cnf *Config) {
	reportIntervalSec := flag.Int64("r", int64(defaultConfig.ReportInterval.Seconds()), "частота отправки метрик на сервер (по умолчанию 10 секунд)")
	cnf.ReportInterval = time.Duration(*reportIntervalSec) * time.Second

	pollIntervalSec := flag.Int64("p", int64(defaultConfig.PollInterval.Seconds()), "частота опроса метрик из пакета runtime (по умолчанию 2 секунды)")
	cnf.PollInterval = time.Duration(*pollIntervalSec) * time.Second

	flag.StringVar(&cnf.Address, "a", defaultConfig.Address, "адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080")

	flag.StringVar(&cnf.Key, "k", defaultConfig.Key, "ключ для подписания запросов SHA256 подписью")

	flag.Int64Var(&cnf.ReportersPoolSize, "l", defaultConfig.ReportersPoolSize, "количество одновременно исходящих запросов на сервер")

	flag.Parse()
}
