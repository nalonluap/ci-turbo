package rules

import (
	"fmt"
	"strings"

	"github.com/nalonluap/ci-turbo/internal/app/service"
	"github.com/nalonluap/ci-turbo/internal/domain"
)

// Убеждаемся, что CachingRule реализует интерфейс Rule.
var _ service.Rule = (*CachingRule)(nil)

var installationCommands = map[string]string{
	"npm install":     "node_modules/",
	"yarn install":    "node_modules/",
	"pip install":     ".venv/",
	"mvn install":     ".m2/",
	"go mod download": "$GOPATH/pkg/mod/",
}

type CachingRule struct{}

func NewCachingRule() *CachingRule {
	return &CachingRule{}
}

func (r *CachingRule) Execute(p *domain.Pipeline) []domain.Finding {
	var findings []domain.Finding

	for _, job := range p.Jobs {
		var detectedCommand string

		for _, scriptLine := range job.Script {
			for command := range installationCommands {
				if strings.Contains(scriptLine, command) {
					detectedCommand = command
					break
				}
			}
			if detectedCommand != "" {
				break
			}
		}

		if detectedCommand == "" {
			continue // В этой задаче нет команд установки, пропускаем
		}

		recommendedPath := installationCommands[detectedCommand]

		// Проверяем, есть ли кэш и содержит ли он рекомендованный путь
		if !containsPath(job.Cache.Paths, recommendedPath) {

			// Формируем сообщение в зависимости от того, есть ли кэш вообще
			message := fmt.Sprintf("Задача выполняет команду '%s', но для нее не настроено кэширование.", detectedCommand)
			suggestion := fmt.Sprintf("Добавьте блок 'cache' в эту задачу с путем к '%s', чтобы значительно ускорить будущие запуски.", recommendedPath)

			if len(job.Cache.Paths) > 0 {
				message = fmt.Sprintf("Задача выполняет команду '%s', но ее кэш не включает необходимую директорию.", detectedCommand)
				suggestion = fmt.Sprintf("Добавьте путь '%s' в секцию cache.paths этой задачи.", recommendedPath)
			}

			finding := domain.Finding{
				RuleID:     "CACHE-002", // Новый ID для более умного правила
				Severity:   "Warning",
				JobName:    job.Name,
				Message:    message,
				Suggestion: suggestion,
			}
			findings = append(findings, finding)
		}
	}

	return findings
}

func containsPath(paths []string, targetPath string) bool {
	for _, p := range paths {
		if p == targetPath {
			return true
		}
	}
	return false
}
