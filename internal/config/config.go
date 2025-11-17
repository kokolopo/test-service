package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	SSLKey     string
}

func LoadConfig() *Config {
	// Hanya memuat .env jika kita tidak berada di lingkungan Docker/Production
	// Jika godotenv.Load() gagal, kita berasumsi variabel sudah dimuat oleh Docker/K8s.
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			// Ini akan log warning, tetapi tidak fatal jika sudah di Docker Compose
			log.Println("Warning: Could not load .env file. Relying on existing environment variables.")
		}
	}

	cfg := &Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "user_service"),
		SSLKey:     getEnv("SSL_KEY", "/pkg/database/ssl.pem"),
	}

	log.Printf("Loaded config: %+v", cfg)
	return cfg
}

// 🚀 getEnv yang diperbaiki: HANYA mengambil variabel, tanpa mencoba memuat .env
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
