package config

import (
	"os"
	"strconv"
)

func LoadEnv() {
	// Server
	if os.Getenv("APP_PORT") == "" {
		os.Setenv("APP_PORT", "8080")
	}
	if os.Getenv("APP_ENV") == "" {
		os.Setenv("APP_ENV", "development")
	}

	// Database
	if os.Getenv("DB_HOST") == "" {
		os.Setenv("DB_HOST", "localhost")
	}
	if os.Getenv("DB_PORT") == "" {
		os.Setenv("DB_PORT", "5432")
	}
	if os.Getenv("DB_USER") == "" {
		os.Setenv("DB_USER", "postgres")
	}
	if os.Getenv("DB_PASSWORD") == "" {
		os.Setenv("DB_PASSWORD", "password")
	}
	if os.Getenv("DB_NAME") == "" {
		os.Setenv("DB_NAME", "admin_panel")
	}

	// Redis
	if os.Getenv("REDIS_HOST") == "" {
		os.Setenv("REDIS_HOST", "localhost")
	}
	if os.Getenv("REDIS_PORT") == "" {
		os.Setenv("REDIS_PORT", "6379")
	}

	// Cloudinary
	if os.Getenv("CLOUDINARY_CLOUD_NAME") == "" {
		os.Setenv("CLOUDINARY_CLOUD_NAME", "demo")
	}
	if os.Getenv("CLOUDINARY_API_KEY") == "" {
		os.Setenv("CLOUDINARY_API_KEY", "demo_key")
	}
	if os.Getenv("CLOUDINARY_API_SECRET") == "" {
		os.Setenv("CLOUDINARY_API_SECRET", "demo_secret")
	}

	// JWT
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "super-secret-key-change-in-prod")
	}
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		initVal, err := strconv.Atoi(value)
		if err == nil {
			return initVal
		}
	}
	return defaultValue
}