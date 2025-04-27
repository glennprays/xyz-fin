package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName                        string
	AppPort                        string
	DBHost                         string
	DBPort                         string
	DBUser                         string
	DBPassword                     string
	DBName                         string
	JWTAccessSecret                string
	JWTRefreshSecret               string
	JWTAccessTokenDurationMinutes  time.Duration
	JWTRefreshTokenDurationMinutes time.Duration
}

func LoadConfig() *Config {
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found (using system envs)")
		}
	}

	jwtAccessTokenDurationStr := getEnv("JWT_ACCESS_TOKEN_DURATION_MINUTES", "5")
	jwtRefreshTokenDurationStr := getEnv("JWT_REFRESH_TOKEN_DURATION_MINUTES", "60")

	jwtAccessTokenDuration, err := strconv.Atoi(jwtAccessTokenDurationStr)
	if err != nil {
		log.Fatal("invalid JWT_ACCESS_TOKEN_DURATION_MINUTES:", err)
	}
	jwtRefreshTokenDuration, err := strconv.Atoi(jwtRefreshTokenDurationStr)
	if err != nil {
		log.Fatal("invalid JWT_REFRESH_TOKEN_DURATION_MINUTES:", err)
	}

	jwtAccessTokenDurationTime := time.Duration(jwtAccessTokenDuration) * time.Minute
	jwtRefreshTokenDurationTime := time.Duration(jwtRefreshTokenDuration) * time.Minute

	return &Config{
		AppName:                        getEnv("APP_NAME", "ordent-test"),
		AppPort:                        getEnv("APP_PORT", "8080"),
		DBHost:                         getEnv("DB_HOST", "localhost"),
		DBPort:                         getEnv("DB_PORT", "5432"),
		DBUser:                         getEnv("DB_USER", "postgres"),
		DBPassword:                     getEnv("DB_PASSWORD", "secret"),
		DBName:                         getEnv("DB_NAME", "mydb"),
		JWTAccessSecret:                getEnv("JWT_ACCESS_SECRET", "supersecret"),
		JWTRefreshSecret:               getEnv("JWT_REFRESH_SECRET", "supersecret"),
		JWTAccessTokenDurationMinutes:  jwtAccessTokenDurationTime,
		JWTRefreshTokenDurationMinutes: jwtRefreshTokenDurationTime,
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
