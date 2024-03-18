package app

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	HTTPAddress      string        `yaml:"HTTP_ADDRESS" env:"HTTP_ADDRESS"`
	HTTPReadTimeout  time.Duration `yaml:"HTTP_READ_TIMEOUT" env:"PG_HOST"`
	HTTPWriteTimeout time.Duration `yaml:"HTTP_WRITE_TIMEOUT" env:"PG_HOST"`
	HTTPIdleTimeout  time.Duration `yaml:"HTTP_IDLE_TIMEOUT" env:"PG_HOST"`

	PgHost      string `yaml:"PG_HOST" env:"PG_HOST"`
	PG_PASSWORD string `yaml:"PG_PASSWORD" env:"PG_PASSWORD"`
	PG_PORT     string `yaml:"PG_PORT" env:"PG_PORT"`
	PG_DBNAME   string `yaml:"PG_DBNAME" env:"PG_DBNAME"`
	PG_USER     string `yaml:"PG_USER" env:"PG_USER"`
	PG_SSLMODE  string `yaml:"PG_SSLMODE" env:"PG_SSLMODE"`
}

func NewConfig(config_path string) (Config, error) {
	cfg := Config{}
	err := cleanenv.ReadConfig(config_path, &cfg)

	if err != nil {
		return Config{}, err
	}
	logrus.Info(cfg)
	return cfg, nil
}
