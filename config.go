package main

import (
	"github.com/BurntSushi/toml"
	"time"
)

type Config struct {
	ApiKey   string        `toml:"api_key"`
	CacheTTL time.Duration `toml:"cache_ttl"`
	City     string
}

func LoadConfig(filename string) (*Config, error) {
	config := &Config{}
	if _, err := toml.DecodeFile(filename, config); err != nil {
		return nil, err
	}

	return config, nil
}

func setConfigDefaults(cfg *Config) {
	if cfg.CacheTTL == 0 {
		cfg.CacheTTL = DefaultCacheTTL
	}
}
