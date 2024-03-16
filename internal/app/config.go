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
}

func NewConfig(config_path string) (Config, error) {
	cfg := Config{}
	err := cleanenv.ReadConfig(config_path, &cfg)

	if err != nil {
		os.Exit(1)
	}

	return cfg, nil
}
