package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App      App
		HTTP     HTTP
		Log      Log
		Postgres Postgres
		Redis    Redis
	}

	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	HTTP struct {
		Port string `env:"HTTP_PORT,required" env-default:"8001"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL,required" env-default:"INFO"`
	}

	Postgres struct {
		Url         string `env:"POSTGRES_URL,required"`
		MigratePath string `env:"POSTGRES_MIGRATE_PATH,required" env-default:"file://migrations"`
	}

	Redis struct {
		Addr     string `env:"REDIS_ADDR,required" env-default:"localhost:6379" example:"localhost:6379"`
		Password string `env:"REDIS_PASSWORD,required"`
		DB       int    `env:"REDIS_DB,required" env-default:"0"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	var err error

	path := os.Getenv("PATH_DOTENV")
	if path != "" {
		if err = cleanenv.ReadConfig(path, cfg); err != nil {
			return nil, err
		}
		return cfg, nil
	}

	if err = cleanenv.ReadEnv(cfg); err == nil {
		return cfg, nil
	}

	return nil, err
}
