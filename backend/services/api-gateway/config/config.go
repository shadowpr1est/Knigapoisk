package config

import (
	"log"
	"os"
)

type Config struct {
	HTTPPort string

	AuthServiceAddr    string
	BookServiceAddr    string
	FileServiceAddr    string
	ReadingServiceAddr string
	ReviewServiceAddr  string
}

func Load() *Config {
	return &Config{
		HTTPPort:           getEnv("HTTP_PORT", ":8080"),
		AuthServiceAddr:    mustEnv("AUTH_SERVICE_ADDR"),
		BookServiceAddr:    mustEnv("BOOK_SERVICE_ADDR"),
		FileServiceAddr:    mustEnv("FILE_SERVICE_ADDR"),
		ReadingServiceAddr: mustEnv("READING_SERVICE_ADDR"),
		ReviewServiceAddr:  mustEnv("REVIEW_SERVICE_ADDR"),
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

