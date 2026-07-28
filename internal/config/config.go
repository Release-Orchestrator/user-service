package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	DBHost   string
	DBPort   string
	DBUser   string
	DBPass   string
	DBName   string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:     getEnv("PORT", "8080"),
		DBHost:   getEnv("DB_HOST", "localhost"),
		DBPort:   getEnv("DB_PORT", "5432"),
		DBUser:   getEnv("DB_USER", "postgres"),
		DBPass:   getEnv("DB_PASS", "postgres"),
		DBName:   getEnv("DB_NAME", "user_db"),
	}
}

func (c *Config) DSN() string {
	return "postgres://" + c.DBUser + ":" + c.DBPass + "@" + c.DBHost + ":" + c.DBPort + "/" + c.DBName + "?sslmode=disable"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}