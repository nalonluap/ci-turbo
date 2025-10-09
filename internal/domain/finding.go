package domain

// Finding представляет собой одну найденную проблему в пайплайне.
type Finding struct {
	RuleID     string // ID правила, которое нашло проблему (напр., "CACHE-001")
	Severity   string // Важность: "Critical", "Warning"
	Message    string // Описание проблемы
	Suggestion string // Предложение по исправлению
	JobName    string // Имя задачи, где найдена проблема
}
