package config

import (
	"github.com/caarlos0/env/v6"
	"log"
)

type environs struct {
	Address string `env:"ADDRESS"`
}

func readEnv() environs {
	envs := environs{}

	err := env.Parse(&envs)
	if err != nil {
		log.Fatalf("failed to parse env variables: %v", err)
	}

	return envs
}

func (e environs) configSetters() []configSetter {
	return []configSetter{
		func(cfg *config) {
			if e.Address != "" {
				log.Printf("overriding address with env variable: %s", e.Address)
				cfg.address = e.Address
			}
		},
	}
}
