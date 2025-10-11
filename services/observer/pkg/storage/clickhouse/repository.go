package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/nalonluap/ci-turbo/internal/domain"
	"github.com/nalonluap/ci-turbo/services/observer/internal/interfaces"
)

// ensure Repository implements the MetricReader interface
var _ interfaces.Reader = (*Repository)(nil)

// Repository encapsulates the connection to a ClickHouse database.
type Repository struct {
	db *sql.DB
}

// NewRepository creates and verifies a new connection to ClickHouse.
func NewRepository(dsn string) (*Repository, error) {
	db := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{dsn},
	})

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to ClickHouse: %w", err)
	}

	return &Repository{db: db}, nil
}

// SlowestJobs executes an analytical SQL query against ClickHouse.
func (r *Repository) SlowestJobs(ctx context.Context, projectID string, period time.Duration) ([]domain.SlowestJobAnalytics, error) {
	var results []domain.SlowestJobAnalytics
	since := time.Now().Add(-period)

	// This complex query calculates all the necessary analytics in one go.
	query := `
        SELECT
            JobName,
            avg(DurationMs) AS avg_duration,
            quantile(0.95)(DurationMs) AS p95_duration,
            count() AS runs_count,
            avg(if(ExitCode = 0, 1, 0)) AS success_rate
        FROM metrics
        WHERE
            ProjectID = ? AND
            StartTime >= ?
        GROUP BY JobName
        ORDER BY avg_duration DESC
        LIMIT 10
    `

	rows, err := r.db.QueryContext(ctx, query, projectID, since)
	if err != nil {
		return nil, fmt.Errorf("failed to execute analytics query: %w", err)
	}
	defer rows.Close()

	// Scan the results into our analytics domain model.
	for rows.Next() {
		var res domain.SlowestJobAnalytics
		if err = rows.Scan(
			&res.JobName,
			&res.AverageDurationMs,
			&res.P95DurationMs,
			&res.RunsCount,
			&res.SuccessRate,
		); err != nil {
			return nil, fmt.Errorf("failed to scan query result row: %w", err)
		}
		results = append(results, res)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return results, nil
}
