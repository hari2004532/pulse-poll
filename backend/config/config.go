package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	MongoURI      string
	MongoDatabase string
	RedisAddr     string
	RedisPassword string
	JWTSecret     string
	FrontendURL   string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Port:          getEnv("PORT", "8080"),
		MongoURI:      getEnv("MONGODB_URI", ""),
		MongoDatabase: getEnv("MONGODB_DATABASE", "pulsepoll"),
		RedisAddr:     getEnv("REDIS_ADDR", ""),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		FrontendURL:   getEnv("FRONTEND_URL", "http://localhost:5173"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
