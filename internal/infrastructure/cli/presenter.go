package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/nalonluap/ci-turbo/internal/domain"
)

// Render - это основная функция Presenter'а.
// Она принимает срез "находок" и выводит их в консоль.
func Render(findings []domain.Finding) {
	if len(findings) == 0 {
		color.Green("✅ Анализ завершен. Проблем не найдено!")
		return
	}

	// --- Печатаем заголовок ---
	header := color.New(color.Bold, color.FgWhite)
	header.Printf("🔍 Анализ завершен. Найдено проблем: %d\n\n", len(findings))

	// --- Создаем "принтеры" для разных уровней серьезности ---
	warningColor := color.New(color.FgYellow)
	jobNameColor := color.New(color.FgCyan, color.Bold)
	suggestionColor := color.New(color.FgWhite).Add(color.Italic)

	// --- Итерируемся по находкам и печатаем их ---
	for i, f := range findings {
		warningColor.Printf("⚠️ [%s] в задаче ", f.Severity)
		jobNameColor.Printf("'%s'\n", f.JobName)

		fmt.Printf("   - Проблема: %s\n", f.Message)

		suggestionColor.Printf("   - Рекомендация: %s\n", f.Suggestion)

		if i < len(findings)-1 {
			fmt.Println() // Добавляем отступ между проблемами
		}
	}
}
