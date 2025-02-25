package config

import (
	"flag"
	"time"
)

type flagsconf struct {
	ConfigFilePath        string
	CryptoKey             string
	Address               string
	Key                   string
	PollIntervalSeconds   int64
	ReportIntervalSeconds int64
	ReportersPoolSize     int64
}

func parseFlags() flagsconf {
	flgs := flagsconf{}

	flag.Int64Var(&flgs.ReportIntervalSeconds, "r", int64(defaultConfig.ReportInterval.Seconds()), "частота отправки метрик на сервер (по умолчанию 10 секунд)")

	flag.Int64Var(&flgs.PollIntervalSeconds, "p", int64(defaultConfig.PollInterval.Seconds()), "частота опроса метрик из пакета runtime (по умолчанию 2 секунды)")

	flag.StringVar(&flgs.Address, "a", defaultConfig.Address, "адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080")

	flag.StringVar(&flgs.Key, "k", defaultConfig.Key, "ключ для подписания запросов SHA256 подписью")

	flag.Int64Var(&flgs.ReportersPoolSize, "l", defaultConfig.ReportersPoolSize, "количество одновременно исходящих запросов на сервер")

	flag.StringVar(&flgs.CryptoKey, "crypto-key", defaultConfig.CryptoKey, "путь до файла с публичным ключом")

	flag.StringVar(&flgs.ConfigFilePath, "c", "", "Path to config file")

	flag.Parse()

	return flgs
}

func (f flagsconf) override(cfg *Config) {
	cfg.overrideAddressIfNotDefault(f.Address, "flags")
	cfg.overrideReportIntervalIfNotDefault(time.Duration(f.ReportIntervalSeconds)*time.Second, "flags")
	cfg.overridePollIntervalIfNotDefault(time.Duration(f.PollIntervalSeconds)*time.Second, "flags")
	cfg.overrideCryptoKeyIfNotDefault(f.CryptoKey, "flags")
	cfg.overrideKeyIfNotDefault(f.Key, "flags")
	cfg.overridePoolSizeIfNotDefault(f.ReportersPoolSize, "flags")
}
