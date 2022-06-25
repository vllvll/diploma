package config

import (
	"github.com/caarlos0/env/v6"
	flag "github.com/spf13/pflag"
)

type GophermartConfig struct {
	Address              string `env:"RUN_ADDRESS"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DatabaseUri          string `env:"DATABASE_URI"`
	Key                  string `env:"KEY"`
}

func CreateGophermartConfig() (*GophermartConfig, error) {
	var cfg GophermartConfig

	flag.StringVarP(&cfg.Address, "address", "a", "127.0.0.1:8080", "Address. Format: ip:port (for example: 127.0.0.1:8080")
	flag.StringVarP(&cfg.DatabaseUri, "database-uri", "d", "", "Database uri. Format: string (for example: postgres://username:password@localhost:5432/database_name)")
	flag.StringVarP(&cfg.AccrualSystemAddress, "accrual-system-address", "r", "", "Accrual system address. Format: ip:port (for example: 127.0.0.1:8080")
	flag.StringVarP(&cfg.Key, "key", "k", "", "Key. Format: string (for example: ?)")

	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
