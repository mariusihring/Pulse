package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string

	Server struct {
		Port string
		Host string
	}

	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SSLMode  string
	}

	JWT struct {
		Secret   string
		Duration time.Duration
	}

	Redis struct {
		Host     string
		Port     string
		Password string
		DB       int
	}

	SendGrid struct {
		APIKey    string
		FromName  string
		FromEmail string
		Templates struct {
			VerifyEmail   string
			ResetPassword string
		}
	}

	ApiKeys struct {
		Blockdaemon string
		Coinbase    string
	}
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Errorf("error loading .env file: %w", err))
	}

	config := &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
	}

	// Server config
	config.Server.Port = getEnv("SERVER_PORT", "8080")
	config.Server.Host = getEnv("SERVER_HOST", "0.0.0.0")

	// Database config
	config.Database.Host = getEnv("DB_HOST", "localhost")
	config.Database.Port = getEnv("DB_PORT", "5432")
	config.Database.User = getEnv("DB_USER", "postgres")
	config.Database.Password = getEnv("DB_PASSWORD", "postgres")
	config.Database.Name = getEnv("DB_NAME", "crypto_tracker")
	config.Database.SSLMode = getEnv("DB_SSLMODE", "disable")

	// JWT config
	config.JWT.Secret = getEnvRequired("JWT_SECRET")
	config.JWT.Duration = time.Duration(getEnvAsInt("JWT_DURATION_HOURS", 24)) * time.Hour

	// Redis config
	config.Redis.Host = getEnv("REDIS_HOST", "localhost")
	config.Redis.Port = getEnv("REDIS_PORT", "6379")
	config.Redis.Password = getEnv("REDIS_PASSWORD", "")
	config.Redis.DB = getEnvAsInt("REDIS_DB", 0)

	// SendGrid config
	config.SendGrid.APIKey = getEnvRequired("SENDGRID_API_KEY")
	config.SendGrid.FromName = getEnv("SENDGRID_FROM_NAME", "Crypto Tracker")
	config.SendGrid.FromEmail = getEnvRequired("SENDGRID_FROM_EMAIL")
	config.SendGrid.Templates.VerifyEmail = getEnvRequired("SENDGRID_TEMPLATE_VERIFY_EMAIL")
	config.SendGrid.Templates.ResetPassword = getEnvRequired("SENDGRID_TEMPLATE_RESET_PASSWORD")

	// API keys
	config.ApiKeys.Blockdaemon = getEnvRequired("BLOCKDAEMON_KEY")
	config.ApiKeys.Coinbase = getEnvRequired("COINBASE_KEY")

	return config
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvRequired(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	panic(fmt.Sprintf("required environment variable %s is not set", key))
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
