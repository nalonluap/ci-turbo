package rules

import (
	"fmt"
	"strings"

	"github.com/nalonluap/ci-turbo/internal/app/service"
	"github.com/nalonluap/ci-turbo/internal/domain"
)

var _ service.Rule = (*ParallelismRule)(nil)

var parallelizableJobKeywords = []string{
	"test",    // Тесты (unit, integration)
	"spec",    // Тесты (BDD)
	"e2e",     // End-to-end тесты
	"qa",      // QA-проверки
	"lint",    // Линтинг
	"build",   // Сборка
	"compile", // Компиляция
	"deploy",  // Развертывание
}

type ParallelismRule struct{}

func NewParallelismRule() *ParallelismRule {
	return &ParallelismRule{}
}

func (r *ParallelismRule) Execute(p *domain.Pipeline) []domain.Finding {
	var findings []domain.Finding

	for _, job := range p.Jobs {
		isCandidate := false

		// 1. Проверяем, является ли задача кандидатом по имени
		for _, keyword := range parallelizableJobKeywords {
			if strings.Contains(job.Name, keyword) {
				isCandidate = true
				break
			}
		}

		if !isCandidate {
			continue
		}

		// 2. Проверяем, что для нее не настроено распараллеливание
		if job.Parallel <= 1 {
			// 3. Создаем "находку" с более общим сообщением
			finding := domain.Finding{
				RuleID:     "PARALLEL-001",
				Severity:   "Warning",
				JobName:    job.Name,
				Message:    fmt.Sprintf("Задача '%s' похожа на длительную операцию (сборка, тестирование, развертывание), но не выполняется параллельно.", job.Name),
				Suggestion: "Если эта задача может быть разделена на независимые части (например, запуск тестов), добавьте `parallel: 5` в ее конфигурацию. Это может значительно сократить время выполнения.",
			}
			findings = append(findings, finding)
		}
	}

	return findings
}
