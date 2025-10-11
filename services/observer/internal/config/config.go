package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

// Config содержит всю конфигурацию для приложения turbo-observer.
type Config struct {
	Server     ServerConfig
	NATS       NATSConfig
	ClickHouse ClickHouseConfig
}

// ServerConfig содержит настройки для самого HTTP-сервера.
type ServerConfig struct {
	Port string `env:"SERVER_PORT" envDefault:"8080"`
}

// NATSConfig содержит настройки для подключения к брокеру NATS.
type NATSConfig struct {
	URL        string `env:"NATS_URL" envDefault:"nats://localhost:4222"`
	StreamName string `env:"NATS_STREAM_NAME" envDefault:"METRICS"`
	Subject    string `env:"NATS_SUBJECT" envDefault:"metrics.*"` // Шаблон для топиков
}

// ClickHouseConfig содержит настройки для подключения к базе данных ClickHouse.
type ClickHouseConfig struct {
	DSN string `env:"CLICKHOUSE_DSN" envDefault:"clickhouse://localhost:9000/default"`
}

// Load читает конфигурацию из переменных окружения и валидирует её.
func Load() (*Config, error) {
	cfg := &Config{}

	// Библиотека env автоматически разбирает переменные окружения в поля структуры.
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("не удалось разобрать конфигурацию приложения: %w", err)
	}

	return cfg, nil
}
