package config

import (
	"os"
	"strconv"
)

func getEnv(key, fallback string) string { 
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	val, err := strconv.Atoi(getEnv(key, ""))
	if err != nil {
		return fallback
	}
	return val
}