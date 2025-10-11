package domain

import "time"

type Event struct {
	ProjectID   string    `json:"project_id"`
	PipelineID  string    `json:"pipeline_id"`
	JobID       string    `json:"job_id"`
	JobName     string    `json:"job_name"`
	RunnerID    string    `json:"runner_id"`
	StartTime   time.Time `json:"start_time"`
	DurationMs  int64     `json:"duration_ms"`
	ExitCode    int       `json:"exit_code"`
	CacheHit    bool      `json:"cache_hit"`     // Пока можно оставить false
	QueueTimeMs int64     `json:"queue_time_ms"` // Время ожидания в очереди
}
