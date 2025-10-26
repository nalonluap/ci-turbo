// services/turbo-observer/internal/service/ingest.go
package service

import (
	"context"

	"github.com/nalonluap/ci-turbo/internal/domain"
	"github.com/nalonluap/ci-turbo/services/observer/internal/interfaces"
)

type IngestService struct {
	broker interfaces.Broker
}

func NewIngestService(broker interfaces.Broker) *IngestService {
	return &IngestService{broker: broker}
}

// Ingest принимает событие и отправляет его в очередь через брокера.
func (s *IngestService) Ingest(ctx context.Context, event *domain.Event) error {
	// Здесь может быть дополнительная логика: обогащение данных, проверка и т.д.

	// Делегируем отправку нашему брокеру (мы не знаем, Kafka это или что-то еще)
	return s.broker.Publish(ctx, "metrics.data", event)
}
