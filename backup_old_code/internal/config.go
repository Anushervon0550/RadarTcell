package config

import (
	"fmt"
	"os"
)

type Config struct {
	Env         string
	AppPort     string
	DatabaseURL string
	RedisAddr   string
}

func Load() (Config, error) {
	cfg := Config{
		Env:       getenv("ENV", "local"),
		AppPort:   getenv("APP_PORT", "8080"),
		RedisAddr: getenv("REDIS_ADDR", "127.0.0.1:6379"),
	}

	cfg.DatabaseURL = os.Getenv("DATABASE_URL")
	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
