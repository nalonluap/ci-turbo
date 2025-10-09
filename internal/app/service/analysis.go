package service

import "github.com/nalonluap/ci-turbo/internal/domain"

// AnalysisService — это ядро бизнес-логики.
// Он оркестрирует процесс анализа, используя реализации интерфейсов,
// которые ему передали при создании. Он не зависит от конкретных технологий.
type AnalysisService struct {
	rules []Rule
}

// NewAnalysisService создаёт новый экземпляр сервиса анализа.
// В качестве зависимости он принимает срез из любых структур,
// которые удовлетворяют интерфейсу Rule.
func NewAnalysisService(rules []Rule) *AnalysisService {
	return &AnalysisService{
		rules: rules,
	}
}

// Analyze выполняет полный анализ пайплайна.
// Он последовательно применяет каждое правило к предоставленной
// модели пайплайна и собирает все найденные проблемы в единый срез.
func (s *AnalysisService) Analyze(p *domain.Pipeline) []domain.Finding {
	var allFindings []domain.Finding

	// Проходим по всем правилам, зарегистрированным в сервисе
	for _, rule := range s.rules {
		// Выполняем правило и получаем срез "находок"
		findings := rule.Execute(p)
		// Добавляем результат в общий список
		allFindings = append(allFindings, findings...)
	}

	return allFindings
}
