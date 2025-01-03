package config

import (
	"database/sql"
	"os"
)

type Config struct {
	Port        string
	DB          *sql.DB
	GithubToken string
}

func New() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		GithubToken: getEnv("GITHUB_TOKEN", ""),
		DB:          nil, // We'll implement DB connection later
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
