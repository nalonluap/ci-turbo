package domain

// Pipeline - это универсальное представление CI/CD пайплайна.
type Pipeline struct {
	Jobs map[string]Job
}

// Job - это универсальное представление задачи в пайплайне.
type Job struct {
	Name   string
	Stage  string
	Script []string
	Cache  CacheConfig
}

// CacheConfig - универсальное представление настроек кэша.
type CacheConfig struct {
	Key   string
	Paths []string
}
