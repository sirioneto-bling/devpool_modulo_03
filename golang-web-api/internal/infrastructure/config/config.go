package config

import (
	"os"
	"strconv"
)

// Config holds all the configuration values loaded from environment variables.
// In a production chassi these would include dozens of fields (Kafka, Redis, Mongo, tracing, etc.).
// Here we keep only what is needed: the app itself and one database.
type Config struct {
	AppName      string
	Port         string
	Environment  string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	DBMaxOpenConns int
	DBMaxIdleConns int
}

// LoadConfig reads environment variables and returns a Config with sensible defaults.
// Every field uses a helper so missing vars never panic -- they just fall back to the default.
func LoadConfig() *Config {
	return &Config{
		AppName:      getEnv("APP_NAME", "devpool-base-web-api"),
		Port:         getEnv("API_PORT", "8080"),
		Environment:  getEnv("ENV", "development"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "3306"),
		DBUser:       getEnv("DB_USER", "devpool"),
		DBPassword:   getEnv("DB_PASSWORD", "devpool123"),
		DBName:       getEnv("DB_NAME", "devpool"),
		DBMaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 10),
		DBMaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 5),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}
