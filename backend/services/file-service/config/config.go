package config

import (
	"log"
	"os"
)

type Config struct {
	GRPCPort    string
	PostgresDSN string

	MinIOEndpoint  string
	MinIOAccessKey string
	MinIOSecretKey string
	MinIOBucket    string
	MinIOUseSSL    bool
}

func Load() *Config {
	return &Config{
		GRPCPort:       getEnv("GRPC_PORT", "0.0.0.0:9093"),
		PostgresDSN:    mustEnv("POSTGRES_DSN"),
		MinIOEndpoint:  mustEnv("MINIO_ENDPOINT"),
		MinIOAccessKey: mustEnv("MINIO_ACCESS_KEY"),
		MinIOSecretKey: mustEnv("MINIO_SECRET_KEY"),
		MinIOBucket:    mustEnv("MINIO_BUCKET"),
		MinIOUseSSL:    getEnv("MINIO_USE_SSL", "false") == "true",
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

