package config

import "time"

type Config struct {
	PollInterval      time.Duration
	ReportInterval    time.Duration
	Address           string
	Key               string
	ReportersPoolSize int64
}

var defaultConfig = Config{
	PollInterval:      2 * time.Second,
	ReportInterval:    10 * time.Second,
	Address:           "localhost:8080",
	Key:               "",
	ReportersPoolSize: 100,
}

func Init() Config {

	cnf := defaultConfig

	parseFlags(&cnf)
	parseEnv(&cnf)

	return cnf
}
