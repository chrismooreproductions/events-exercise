package config

import (
	"os"
)

type Config struct {
	Host string
	Port string
}

func LoadConfig() Config {
	return Config{
		Host: getEnv("HOST", "localhost"),
		Port: getEnv("PORT", "9001"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
