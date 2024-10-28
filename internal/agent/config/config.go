package config

import "time"

type Config struct {
	PollInterval   time.Duration
	ReportInterval time.Duration
	Address        string
	Key            string
}

var defaultConfig = Config{
	PollInterval:   2 * time.Second,
	ReportInterval: 10 * time.Second,
	Address:        "localhost:8080",
	Key:            "",
}

func Init() Config {

	cnf := defaultConfig

	parseFlags(&cnf)
	parseEnv(&cnf)

	return cnf
}
