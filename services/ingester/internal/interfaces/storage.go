package interfaces

import (
	"context"

	"github.com/nalonluap/ci-turbo/internal/domain"
)

type Writer interface {
	Batch(ctx context.Context, events []*domain.Event) error
}
