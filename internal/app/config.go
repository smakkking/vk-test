package app

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
}

func NewConfig(config_path string) (*Config, error) {
	cfg := Config{}
	err := cleanenv.ReadConfig(config_path, &cfg)

	if err != nil {
		os.Exit(1)
	}

	return &cfg, nil
}
