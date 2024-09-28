package config

import (
	"flag"
	"log"
)

type flags struct {
	Address string
}

func parseFlags() flags {
	flgs := flags{}
	flag.StringVar(&flgs.Address, "a", "", "Address of server, default: localhost:8080")
	flag.Parse()
	return flgs
}

func (f flags) configSetters() []configSetter {
	return []configSetter{
		func(cfg *config) {
			if f.Address != "" {
				log.Printf("overriding address with flag: %s", f.Address)
				cfg.address = f.Address
			}
		},
	}
}
