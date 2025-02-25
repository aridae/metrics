package config

import "time"

type Config struct {
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

func Obtain() Config {

	cnf := defaultConfig

	parseFlags(&cnf)
	parseEnv(&cnf)

	return cnf
}
