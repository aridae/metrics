package config

import "flag"

type Flags struct {
	address string
}

func parseFlags() Flags {
	flags := Flags{}
	flag.StringVar(&flags.address, "a", "localhost:8080", "Address of server, default: localhost:8080")
	flag.Parse()
	return flags
}

func (f Flags) configSetters() []configSetter {
	return []configSetter{
		func(cfg *config) { cfg.address = f.address },
	}
}
