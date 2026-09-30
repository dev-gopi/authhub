package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	RabbitMQ RabbitMQConfig
	Vault    VaultConfig
}

type AppConfig struct {
	Environment     string
	AuthAPIAddr     string
	ControlAPIAddr  string
	RequestTimeout  time.Duration
	ShutdownTimeout time.Duration
	IssuerBaseURL   string
}

type PostgresConfig struct {
	DSN string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type RabbitMQConfig struct {
	URL string
}

type VaultConfig struct {
	Address string
	Token   string

	TransitMount string
}

func Load() (*Config, error) {
	// Local development convenience.
	// In production environment variables should be injected.
	_ = godotenv.Load()

	requestTimeout, err := parseDuration("REQUEST_TIMEOUT", "30s")
	if err != nil {
		return nil, err
	}

	shutdownTimeout, err := parseDuration("SHUTDOWN_TIMEOUT", "15s")
	if err != nil {
		return nil, err
	}

	redisDB, err := parseInt("REDIS_DB", 0)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		App: AppConfig{
			Environment:     getEnv("APP_ENV", "development"),
			AuthAPIAddr:     getEnv("AUTH_API_ADDR", ":8080"),
			ControlAPIAddr:  getEnv("CONTROL_API_ADDR", ":8081"),
			RequestTimeout:  requestTimeout,
			ShutdownTimeout: shutdownTimeout,
			IssuerBaseURL: getEnv(
				"ISSUER_BASE_URL",
				"http://localhost:8080",
			),
		},

		Postgres: PostgresConfig{
			DSN: os.Getenv("POSTGRES_DSN"),
		},

		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       redisDB,
		},

		RabbitMQ: RabbitMQConfig{
			URL: os.Getenv("RABBITMQ_URL"),
		},

		Vault: VaultConfig{
			Address: os.Getenv("VAULT_ADDR"),
			Token:   os.Getenv("VAULT_TOKEN"),

			TransitMount: getEnv(
				"VAULT_TRANSIT_MOUNT",
				"transit",
			),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.Postgres.DSN == "" {
		return fmt.Errorf("POSTGRES_DSN is required")
	}

	if c.Redis.Addr == "" {
		return fmt.Errorf("REDIS_ADDR is required")
	}

	if c.RabbitMQ.URL == "" {
		return fmt.Errorf("RABBITMQ_URL is required")
	}

	if c.Vault.Address == "" {
		return fmt.Errorf("VAULT_ADDR is required")
	}

	return nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func parseDuration(key, fallback string) (time.Duration, error) {
	value := getEnv(key, fallback)

	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}

	return duration, nil
}

func parseInt(key string, fallback int) (int, error) {
	value := os.Getenv(key)

	if value == "" {
		return fallback, nil
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}

	return n, nil
}
