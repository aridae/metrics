package config

import (
	"flag"
	"time"
)

type flagsconf struct {
	ConfigFilePath          string
	CryptoKey               string
	AddressOverride         string
	FileStoragePathOverride string
	DatabaseDsnOverride     string
	Key                     string
	StoreIntervalOverride   time.Duration
	RestoreOverride         bool
}

func parseFlags() flagsconf {
	flgs := flagsconf{}

	flag.StringVar(&flgs.AddressOverride, "a", addressDefaultVal, "Address of server")

	storeInterval := flag.Int64("i", int64(storeIntervalDefault.Seconds()), "Backup store interval")
	flgs.StoreIntervalOverride = time.Duration(*storeInterval) * time.Second

	flag.StringVar(&flgs.FileStoragePathOverride, "f", fileStoragePathDefault, "Backup file path")

	flag.BoolVar(&flgs.RestoreOverride, "r", restoreDefault, "Restore from backup file on start")

	flag.StringVar(&flgs.DatabaseDsnOverride, "d", "", "Database DSN")

	flag.StringVar(&flgs.Key, "k", "", "ключ для подписания запросов SHA256 подписью")

	flag.StringVar(&flgs.CryptoKey, "crypto-key", "", "путь до файла с приватным ключом")

	flag.StringVar(&flgs.ConfigFilePath, "c", "", "Path to config file")

	flag.Parse()

	return flgs
}

func (f flagsconf) override(cfg *Config) {
	cfg.overrideAddressIfNotDefault(f.AddressOverride, "flags")
	cfg.overrideStoreIntervalIfNotDefault(f.StoreIntervalOverride, "flags")
	cfg.overrideFileStoragePathIfNotDefault(f.FileStoragePathOverride, "flags")
	cfg.overrideRestoreIfNotDefault(f.RestoreOverride, "flags")
	cfg.overrideCryptoKeyIfNotDefault(f.CryptoKey, "flags")
}
