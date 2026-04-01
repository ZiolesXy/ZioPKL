package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	Environment string
	DB DBConfig
	Redis RedisConfig
	JwtSecret string
	MasterKey string
	FCMKey string
	CorsOrigins []string
}

type DBConfig struct {
	Host string
	Port string
	User string
	Password string
	Name string
}

type RedisConfig struct {
	Host string
	Port string
	Password string
}

func LoadEnv() *Config {
	var c Config

	if err := godotenv.Load(); err != nil {
		log.Panicln("No .env file found, using OS defaults")
	}

	c.Environment = getEnv("ENV", "development")
	c.Port = getEnv("SERVER_PORT", "8080")
	c.JwtSecret = getEnv("JWT_SECRET", "change_me")
	c.MasterKey = getEnv("E2EE_MASTER_KEY", "")

	c.DB.Host = getEnv("DB_HOST", "localhost")
	c.DB.Port = getEnv("DB_PORT", "5432")
	c.DB.User = getEnv("DB_USER", "postgres")
	c.DB.Password = getEnv("DB_PASSWORD", "password")
	c.DB.Name = getEnv("DB_NAME", "chatapp_db")

	c.Redis.Host = getEnv("REDIS_HOST", "localhost")
	c.Redis.Port = getEnv("REDIS_PORT", "6379")
	c.Redis.Password = getEnv("REDIS_PASSWORD", "")

	c.FCMKey = getEnv("FCM_SERVER_KEY", "")

	c.CorsOrigins = []string{getEnv("ALLOWED_ORIGINS", "*")}

	return &c
}

func (c *Config) DSN() string {
	return "host=" + c.DB.Host +
		" port=" + c.DB.Port +
		" user=" + c.DB.User +
		" password=" + c.DB.Password +
		" dbname=" + c.DB.Name +
		" sslmode=disable"
}

func (c *Config) RedisURL() string {
	return  c.Redis.Host + ":" + c.Redis.Port
}