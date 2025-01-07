package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	Env *Config
)

type Config struct {
	// Server configs
	Host        string
	Port        string
	Environment string

	// Database configs
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPoolSize string

	// CORS configs
	CORSOrigins []string
	CORSMethods []string
	CORSHeaders []string
	CORSMaxAge  int

	//Google Drive configs
	GOOGLE_DRIVE_CREDENTIALS_PATH string
	GOOGLE_DRIVE_TOKEN_PATH       string
	GOOGLE_DRIVE_REDIRECT_URL     string

	//Kafka configs
	KafkaBrokers string
	KafkaGroupID string
}

func Load() error {
	if err := godotenv.Load(".env"); err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	Env = &Config{
		// Server configs
		Host:        getEnvOrDefault("HOST", "localhost"),
		Port:        getEnvOrDefault("PORT", "8080"),
		Environment: getEnvOrDefault("GO_ENV", "DEV"),

		// Database configs
		DBHost:     getEnvOrDefault("DB_HOST", "localhost"),
		DBPort:     getEnvOrDefault("DB_PORT", "27017"),
		DBUser:     getEnvOrDefault("DB_USER", "mongodb"),
		DBPassword: getEnvOrDefault("DB_PASSWORD", "mongodb"),
		DBName:     getEnvOrDefault("DB_NAME", "mongodb"),
		DBPoolSize: getEnvOrDefault("DB_POOL_SIZE", "100"),

		// CORS configs
		CORSOrigins: []string{getEnvOrDefault("CORS_ORIGINS", "http://localhost:3000")},
		CORSMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		CORSHeaders: []string{"Content-Type", "Authorization"},
		CORSMaxAge:  3600,

		//Google Drive configs
		GOOGLE_DRIVE_CREDENTIALS_PATH: getEnvOrDefault("GOOGLE_DRIVE_CREDENTIALS_PATH", ""),
		GOOGLE_DRIVE_TOKEN_PATH:       getEnvOrDefault("GOOGLE_DRIVE_TOKEN_PATH", ""),
		GOOGLE_DRIVE_REDIRECT_URL:     getEnvOrDefault("GOOGLE_DRIVE_REDIRECT_URL", ""),

		//Kafka configs
		KafkaBrokers: getEnvOrDefault("KAFKA_BROKERS", "172.29.96.1:9092"),
		KafkaGroupID: getEnvOrDefault("KAFKA_GROUP_ID", "my-group"),
	}

	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
