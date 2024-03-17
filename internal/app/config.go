package app

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Address          string
	HTTPReadTimeout  time.Duration
	HTTPWriteTimeout time.Duration
	HTTPIdleTimeout  time.Duration

	PG_HOST     string
	PG_PASSWORD string
	PG_PORT     string
	PG_DBNAME   string
	PG_USER     string
	PG_SSLMODE  string
}

func NewConfig(config_path string) (Config, error) {
	cfg := Config{}
	err := cleanenv.ReadConfig(config_path, &cfg)

	if err != nil {
		os.Exit(1)
	}

	return cfg, nil
}
