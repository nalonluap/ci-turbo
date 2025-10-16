package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

// Config holds all configuration for the ingestion-worker application.
type Config struct {
	NATS       NATSConfig
	ClickHouse ClickHouseConfig
}

// NATSConfig holds NATS-specific configuration.
type NATSConfig struct {
	URL           string `env:"NATS_URL" envDefault:"nats://localhost:4222"`
	StreamName    string `env:"NATS_STREAM_NAME" envDefault:"METRICS"`
	Subject       string `env:"NATS_SUBJECT" envDefault:"metrics.*"`
	ConsumerGroup string `env:"NATS_CONSUMER_GROUP" envDefault:"ingestion-worker-group"`
}

// ClickHouseConfig holds ClickHouse-specific configuration.
type ClickHouseConfig struct {
	Host     string `env:"CLICKHOUSE_HOST" envDefault:"clickhouse"`
	Port     string `env:"CLICKHOUSE_PORT" envDefault:"9000"`
	Database string `env:"CLICKHOUSE_DATABASE" envDefault:"default"`
	User     string `env:"CLICKHOUSE_USER" envDefault:"default"`
	Password string `env:"CLICKHOUSE_PASSWORD"` // Пароль может быть пустым
}

// Load reads configuration from environment variables and validates it.
func Load() (*Config, error) {
	cfg := &Config{}

	// The env library automatically parses variables into the struct fields.
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("could not parse application config: %w", err)
	}

	return cfg, nil
}
