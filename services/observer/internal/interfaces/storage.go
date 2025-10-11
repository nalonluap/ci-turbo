package interfaces

import (
	"context"
	"time"

	"github.com/nalonluap/ci-turbo/internal/domain"
)

// Writer определяет контракт для записи метрик в хранилище.
type Writer interface {
	Batch(ctx context.Context, events []*domain.Event) error
}

// Reader определяет контракт для чтения аналитики.
type Reader interface {
	// SlowestJobs должен вернуть агрегированные данные о самых медленных задачах.
	SlowestJobs(ctx context.Context, projectID string, period time.Duration) ([]domain.SlowestJobAnalytics, error)
}
