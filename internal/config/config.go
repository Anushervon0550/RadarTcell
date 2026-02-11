package config

import (
	"log"
	"os"
)

type Config struct {
	AppPort     string
	DatabaseURL string
	RedisAddr   string
}

func Load() Config {
	cfg := Config{
		AppPort:     getenv("APP_PORT", "8080"),
		DatabaseURL: getenv("DATABASE_URL", ""),
		RedisAddr:   getenv("REDIS_ADDR", "127.0.0.1:6379"),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is empty")
	}

	return cfg
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
