package main

import (
	"fmt"
	"os"

	// Импортируем наши внутренние пакеты
	"github.com/nalonluap/ci-turbo/internal/app/service"
	"github.com/nalonluap/ci-turbo/internal/infrastructure/adapter/gitlab"
	"github.com/nalonluap/ci-turbo/internal/infrastructure/cli"
	"github.com/nalonluap/ci-turbo/internal/infrastructure/rules"
)

func main() {
	// --- Слой Инфраструктуры: Создание конкретных реализаций ---

	// 1. Создаем репозиторий для получения данных (в данном случае, парсер GitLab)
	gitlabParser := gitlab.NewParser()

	// 2. Создаем все правила анализа, которые у нас есть
	cachingRule := rules.NewCachingRule()
	parallelismRule := rules.NewParallelismRule()

	// 3. Создаем сервис анализа и передаем ему все правила
	analysisSvc := service.NewAnalysisService([]service.Rule{
		cachingRule,
		parallelismRule,
	})

	// --- Слой Адаптеров: Создание и запуск CLI ---

	// 4. Создаем корневую команду CLI и внедряем в нее зависимости
	rootCmd := cli.NewRootCommand(analysisSvc, gitlabParser)

	// 5. Запускаем приложение
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка выполнения: %v\n", err)
		os.Exit(1)
	}
}
