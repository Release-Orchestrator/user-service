// Package config provides application configuration.
package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration loaded from environment variables.
type Config struct {
	Port   string
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
}

// Load reads configuration from environment variables.
func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:   getEnv("PORT", "8080"),
		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_PORT", "5432"),
		DBUser: getEnv("DB_USER", "postgres"),
		DBPass: getEnv("DB_PASS", "postgres"),
		DBName: getEnv("DB_NAME", "user_db"),
	}
}

// DSN returns the PostgreSQL connection string.
func (c *Config) DSN() string {
	return "postgres://" + c.DBUser + ":" + c.DBPass + "@" + c.DBHost + ":" + c.DBPort + "/" + c.DBName + "?sslmode=disable"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
