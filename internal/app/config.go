package app

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPAddress      string        `yaml:"HTTP_ADDRESS" env:"HTTP_ADDRESS" env-default:"localhost"`
	HTTPReadTimeout  time.Duration `yaml:"HTTP_READ_TIMEOUT" env:"PG_HOST" env-default:"localhost"`
	HTTPWriteTimeout time.Duration `yaml:"HTTP_WRITE_TIMEOUT" env:"PG_HOST" env-default:"localhost"`
	HTTPIdleTimeout  time.Duration `yaml:"HTTP_IDLE_TIMEOUT" env:"PG_HOST" env-default:"localhost"`

	PgHost      string `yaml:"PG_HOST" env:"PG_HOST" env-default:"localhost"`
	PG_PASSWORD string `yaml:"PG_PASSWORD" env:"PG_PASSWORD" env-default:"localhost"`
	PG_PORT     string `yaml:"PG_PORT" env:"PG_PORT" env-default:"localhost"`
	PG_DBNAME   string `yaml:"PG_DBNAME" env:"PG_DBNAME" env-default:"localhost"`
	PG_USER     string `yaml:"PG_USER" env:"PG_USER" env-default:"localhost"`
	PG_SSLMODE  string `yaml:"PG_SSLMODE" env:"PG_SSLMODE" env-default:"localhost"`
}

func NewConfig(config_path string) (Config, error) {
	cfg := Config{}
	err := cleanenv.ReadConfig(config_path, &cfg)

	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
