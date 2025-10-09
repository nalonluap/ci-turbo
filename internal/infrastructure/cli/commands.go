package cli

import (
	"fmt"

	"github.com/nalonluap/ci-turbo/internal/app/service"
	"github.com/nalonluap/ci-turbo/internal/domain"

	"github.com/spf13/cobra"
)

// NewRootCommand создает корневую команду CLI.
func NewRootCommand(analyzer *service.AnalysisService, repo service.PipelineRepository) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "ci-turbo",
		Short: "CI-Turbo - это умный анализатор CI/CD пайплайнов.",
		Long:  `CI-Turbo помогает находить узкие места в ваших пайплайнах и предлагает способы их ускорения.`,
	}

	// Создаем и добавляем дочернюю команду analyze
	analyzeCmd := createAnalyzeCommand(analyzer, repo)
	rootCmd.AddCommand(analyzeCmd)

	return rootCmd
}

func createAnalyzeCommand(analyzer *service.AnalysisService, repo service.PipelineRepository) *cobra.Command {
	var filePath string

	analyzeCmd := &cobra.Command{
		Use:   "analyze",
		Short: "Анализирует CI/CD конфигурационный файл.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// 1. Получаем доменную модель через репозиторий (адаптер).
			pipeline, err := repo.Get(filePath)

			if err != nil {
				return fmt.Errorf("не удалось получить конфигурацию пайплайна: %w", err)
			}

			// 2. Вызываем сервис анализа с этой моделью.
			findings := analyzer.Analyze(pipeline)

			// 3. Отображаем результат.
			presentFindings(findings)

			return nil
		},
	}

	// Добавляем флаг для указания пути к файлу
	analyzeCmd.Flags().StringVarP(&filePath, "file", "f", ".gitlab-ci.yml", "Путь к конфигурационному файлу CI/CD")

	return analyzeCmd
}

// presentFindings - это адаптер вывода (Presenter).
func presentFindings(findings []domain.Finding) {
	if len(findings) == 0 {
		fmt.Println("✅ Анализ завершен. Проблем не найдено!")
		return
	}

	fmt.Printf("🔍 Анализ завершен. Найдено проблем: %d\n", len(findings))
	// TODO: Реализовать красивый вывод с цветами и таблицами.
	for _, f := range findings {
		fmt.Printf("[%s] %s: %s\n", f.Severity, f.JobName, f.Message)
	}
}
