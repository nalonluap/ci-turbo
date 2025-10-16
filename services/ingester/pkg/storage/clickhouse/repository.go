package clickhouse

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/nalonluap/ci-turbo/internal/domain"
	"github.com/nalonluap/ci-turbo/services/ingester/internal/config"
	"github.com/nalonluap/ci-turbo/services/ingester/internal/interfaces"
)

// ensure Repository implements the MetricWriter interface
var _ interfaces.Writer = (*Repository)(nil)

// Repository encapsulates the connection to a ClickHouse database.
type Repository struct {
	db *sql.DB
}

// NewRepository creates and verifies a new connection to ClickHouse.
func NewRepository(cfg config.ClickHouseConfig) (*Repository, error) {
	db := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)},
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.User,
			Password: cfg.Password,
		},
	})

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("не удалось подключиться к ClickHouse: %w", err)
	}

	return &Repository{db: db}, nil
}

// Batch implements high-performance batch insertion of metrics into ClickHouse.
func (r *Repository) Batch(ctx context.Context, events []*domain.Event) error {
	// 1. Begin a transaction for atomic writes.
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// 2. Prepare the insert statement at once for efficiency.
	stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO metrics (
            ProjectID, PipelineID, JobID, JobName, RunnerID, 
            StartTime, DurationMs, ExitCode, CacheHit, QueueTimeMs
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `)
	if err != nil {
		_ = tx.Rollback() // Attempt to roll back on error
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// 3. Iterate through events and add them to the batch.
	for _, event := range events {
		_, err := stmt.ExecContext(ctx,
			event.ProjectID,
			event.PipelineID,
			event.JobID,
			event.JobName,
			event.RunnerID,
			event.StartTime,
			event.DurationMs,
			event.ExitCode,
			event.CacheHit,
			event.QueueTimeMs,
		)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("error executing statement in batch: %w", err)
		}
	}

	// 4. Commit the transaction to finalize the writing.
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
