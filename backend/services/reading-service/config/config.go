package config

import (
	"log"
	"os"
)

type Config struct {
	GRPCPort    string
	PostgresDSN string
}

func Load() *Config {
	return &Config{
		GRPCPort:    getEnv("GRPC_PORT", "0.0.0.0:9094"),
		PostgresDSN: mustEnv("POSTGRES_DSN"),
	}
}

func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("missing required env %s", key)
	}
	return val
}

func getEnv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

