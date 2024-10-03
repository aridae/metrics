package config

import (
	"flag"
	"log"
	"time"
)

type flags struct {
	AddressOverride              string
	StoreIntervalSecondsOverride int64
	FileStoragePathOverride      string
	RestoreOverride              bool
}

func parseFlags() flags {
	flgs := flags{}

	flag.StringVar(&flgs.AddressOverride, "a", "", "Address of server")
	flag.Int64Var(&flgs.StoreIntervalSecondsOverride, "i", 0, "Address of server")
	flag.StringVar(&flgs.FileStoragePathOverride, "f", "", "Address of server")
	flag.BoolVar(&flgs.RestoreOverride, "r", false, "Address of server")

	flag.Parse()

	return flgs
}

func (f flags) override(cfg *Config) {
	if f.AddressOverride != "" {
		log.Printf("overriding Address with flag: %s", f.AddressOverride)
		cfg.Address = f.AddressOverride
	}

	if f.StoreIntervalSecondsOverride != 0 {
		log.Printf("overriding StoreInterval with flag: %d", f.StoreIntervalSecondsOverride)
		cfg.StoreInterval = time.Duration(f.StoreIntervalSecondsOverride) * time.Second
	}

	if f.FileStoragePathOverride != "" {
		log.Printf("overriding FileStoragePath with flag: %s", f.FileStoragePathOverride)
		cfg.FileStoragePath = f.FileStoragePathOverride
	}

	if f.RestoreOverride {
		log.Printf("overriding Restore with flag: %t", f.RestoreOverride)
		cfg.Restore = f.RestoreOverride
	}
}
