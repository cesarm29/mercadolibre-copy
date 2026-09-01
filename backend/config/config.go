package config

import "os"

type Config struct {
	DatabaseURL string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	JWTSecret   string
	ServerPort  string
}

func Load() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", ""),
		DBHost:      getEnv("PGHOST", getEnv("DB_HOST", "localhost")),
		DBPort:      getEnv("PGPORT", getEnv("DB_PORT", "5432")),
		DBUser:      getEnv("PGUSER", getEnv("DB_USER", "postgres")),
		DBPassword:  getEnv("PGPASSWORD", getEnv("DB_PASSWORD", "")),
		DBName:      getEnv("PGDATABASE", getEnv("DB_NAME", "marketplace")),
		DBSSLMode:   getEnv("DB_SSLMODE", "require"),
		JWTSecret:   getEnv("JWT_SECRET", "super-secret-change-me"),
		ServerPort:  getEnv("PORT", getEnv("SERVER_PORT", "8080")),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
