package config

type Config struct {
	Server ServerConfig
	Database DatabaseConfig
	Redis RedisConfig
	Cloudinary CloudinaryConfig
	JWT JWTConfig
}

type ServerConfig struct {
	Port string
	Env string
}

type DatabaseConfig struct {
	Host string
	Port string
	User string
	Password string
	DBName string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type CloudinaryConfig struct {
	CloudName string
	APIKey    string
	APISecret string
}

type JWTConfig struct {
	Secret string
}

func Load() *Config {
	LoadEnv()

	return &Config{
		Server: ServerConfig{
			Port: GetEnv("APP_PORT", "8080"),
			Env:  GetEnv("APP_ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     GetEnv("DB_HOST", "localhost"),
			Port:     GetEnv("DB_PORT", "5432"),
			User:     GetEnv("DB_USER", "postgres"),
			Password: GetEnv("DB_PASSWORD", "password"),
			DBName:   GetEnv("DB_NAME", "admin_panel"),
		},
		Redis: RedisConfig{
			Host:     GetEnv("REDIS_HOST", "localhost"),
			Port:     GetEnv("REDIS_PORT", "6379"),
			Password: GetEnv("REDIS_PASSWORD", ""),
			DB:       GetEnvInt("REDIS_DB", 0),
		},
		Cloudinary: CloudinaryConfig{
			CloudName: GetEnv("CLOUDINARY_CLOUD_NAME", ""),
			APIKey:    GetEnv("CLOUDINARY_API_KEY", ""),
			APISecret: GetEnv("CLOUDINARY_API_SECRET", ""),
		},
		JWT: JWTConfig{
			Secret: GetEnv("JWT_SECRET", "super-secret-key-change-in-prod"),
		},
	}
}