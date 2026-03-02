package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	GRPCPort    string
	PostgresDSN string

	JWTSecret string
	JWTTTL    time.Duration
}

func Load() *Config {
	return &Config{
		GRPCPort:    getEnv("GRPC_PORT", "0.0.0.0:9091"),
		PostgresDSN: mustEnv("POSTGRES_DSN"),
		JWTSecret:   mustEnv("JWT_SECRET_KEY"),
		JWTTTL:      15 * time.Minute,
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

