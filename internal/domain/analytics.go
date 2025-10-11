package domain

type SlowestJobAnalytics struct {
	JobName           string  `json:"job_name"`
	AverageDurationMs int64   `json:"average_duration_ms"`
	P95DurationMs     int64   `json:"p95_duration_ms"` // 95-й перцентиль времени выполнения
	RunsCount         int     `json:"runs_count"`
	SuccessRate       float64 `json:"success_rate"`
}
