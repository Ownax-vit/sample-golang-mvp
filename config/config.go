package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App     App
		HTTP    HTTP
		Log     Log
		Postgres Postgres
	}

	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	HTTP struct {
		Port           string `env:"HTTP_PORT,required" env-default:"8001"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL,required" end-default:"INFO"`
	}

	Postgres struct {
		Url string `end:"POSRGRES-URL, required"`
	}

)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
