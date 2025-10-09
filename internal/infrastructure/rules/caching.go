package rules

import (
	"fmt"

	"github.com/nalonluap/ci-turbo/internal/app/service"
	"github.com/nalonluap/ci-turbo/internal/domain"
)

// Убеждаемся, что CachingRule реализует интерфейс Rule.
var _ service.Rule = (*CachingRule)(nil)

type CachingRule struct{}

func NewCachingRule() *CachingRule {
	return &CachingRule{}
}

func (r *CachingRule) Execute(p *domain.Pipeline) []domain.Finding {
	// TODO: Реализовать логику поиска проблем с кэшированием.
	// 1. Пройти по всем p.Jobs.
	// 2. В каждой задаче проверить поле Script на наличие команд типа "npm install".
	// 3. Если команда найдена, проверить, что поле Cache у задачи настроено.
	// 4. Если нет - создать и вернуть domain.Finding.
	fmt.Println("INFO: Executing CachingRule...")
	return nil // Placeholder
}
