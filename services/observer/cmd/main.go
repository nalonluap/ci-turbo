package main

import (
	"log"

	// Configuration
	"github.com/nalonluap/ci-turbo/services/observer/internal/config"

	// Application Layer
	"github.com/nalonluap/ci-turbo/services/observer/internal/service"

	// Infrastructure Layer (Entrypoint & Adapters)
	"github.com/nalonluap/ci-turbo/services/observer/pkg/api"
	"github.com/nalonluap/ci-turbo/services/observer/pkg/broker/nats"
	"github.com/nalonluap/ci-turbo/services/observer/pkg/storage/clickhouse"
)

func main() {
	log.Println("Starting Turbo-Observer API server...")

	// --- 1. Load Configuration ---
	// Load configuration from environment variables.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// --- 2. Initialize Infrastructure Adapters ---
	// This is the only place where we have knowledge of concrete technologies (NATS, ClickHouse).

	// Initialize the NATS message broker producer.
	messageBroker, err := nats.NewProducer(cfg.NATS.URL, cfg.NATS.StreamName, cfg.NATS.Subject)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer messageBroker.Close() // Ensure the connection is closed on shutdown.
	log.Println("Connected to NATS broker")

	// Initialize the ClickHouse repository for reading analytics.
	metricStorage, err := clickhouse.NewRepository(cfg.ClickHouse) // Передаем структуру конфига
	if err != nil {
		log.Fatalf("Не удалось подключиться к ClickHouse: %v", err)
	}
	log.Println("Connected to ClickHouse database")

	// --- 3. Initialize Application Services ---
	// Inject the concrete adapter implementations into our services, which depend only on interfaces.
	ingestService := service.NewIngestService(messageBroker)
	queryService := service.NewQueryService(metricStorage)

	// --- 4. Create and Run the API Server ---
	// Pass the application services to the server, which will handle HTTP requests.
	server := api.NewServer(cfg.Server.Port, ingestService, queryService)

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err = server.Run(); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}
