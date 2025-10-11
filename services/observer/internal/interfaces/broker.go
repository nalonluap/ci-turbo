package interfaces

import (
	"context"

	"github.com/nalonluap/ci-turbo/internal/domain"
)

// Broker определяет контракт для отправки сообщений в очередь.
type Broker interface {
	Publish(ctx context.Context, topic string, event *domain.Event) error
}
