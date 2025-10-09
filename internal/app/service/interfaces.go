package service

import (
	"github.com/nalonluap/ci-turbo/internal/domain"
)

// PipelineRepository определяет порт для получения данных о пайплайне.
// Слой приложения не знает, откуда берутся данные: из файла, базы данных или по API.
// Он просто знает, что есть способ получить Pipeline по указанному пути.
type PipelineRepository interface {
	Get(path string) (*domain.Pipeline, error)
}

// Rule определяет порт для одного правила анализа.
// Слой приложения не знает, какова конкретная логика правила, он просто
// доверяет, что любое правило можно выполнить над пайплайном и получить результат.
type Rule interface {
	Execute(p *domain.Pipeline) []domain.Finding
}
