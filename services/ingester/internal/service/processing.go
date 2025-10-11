package service

import (
	"context"
	"fmt"

	"github.com/nalonluap/ci-turbo/internal/domain" // Импортируем общую доменную модель
	"github.com/nalonluap/ci-turbo/services/ingester/internal/interfaces"
)

// ProcessingService инкапсулирует бизнес-логику для обработки и сохранения метрик.
// Он действует как "дирижер" для воркера, координируя валидацию и запись данных.
type ProcessingService struct {
	writer interfaces.Writer // Зависимость от интерфейса, а не от конкретной БД
}

// NewProcessingService создает новый экземпляр ProcessingService.
// В качестве зависимости он принимает любую структуру, реализующую интерфейс MetricWriter.
func NewProcessingService(writer interfaces.Writer) *ProcessingService {
	return &ProcessingService{writer: writer}
}

// ProcessBatch - это основной метод, который вызывается из main-цикла воркера.
// Он применяет бизнес-правила к пачке событий перед их сохранением.
func (s *ProcessingService) ProcessBatch(ctx context.Context, events []*domain.Event) error {
	// --- Шаг 1: Бизнес-логика - Валидация и фильтрация ---
	// Это критически важный этап, который обеспечивает чистоту данных в нашей БД.

	validEvents := make([]*domain.Event, 0, len(events))
	for _, event := range events {
		// Пример правила: событие считается невалидным, если у него нет ID проекта или ID задачи.
		if event.ProjectID == "" || event.JobID == "" {
			// В реальном приложении здесь будет более продвинутое логирование.
			fmt.Printf("ПРЕДУПРЕЖДЕНИЕ: Пропущено невалидное событие без ProjectID или JobID: %+v\n", event)
			continue // Переходим к следующему событию в пачке
		}

		// Здесь можно добавить другую логику в будущем:
		// - Обогащение данных (например, добавить название проекта по его ID).
		// - Проверка на аномальные значения (например, длительность выполнения > 24 часов).

		validEvents = append(validEvents, event)
	}

	// --- Шаг 2: Проверка, есть ли что записывать ---
	// Если после фильтрации не осталось валидных событий, нет смысла обращаться к БД.
	if len(validEvents) == 0 {
		fmt.Println("ИНФО: В пачке не найдено валидных событий для записи.")
		return nil
	}

	// --- Шаг 3: Делегирование записи хранилищу ---
	// Сервис выполнил свою работу (обработал данные) и теперь передает их
	// адаптеру хранилища для фактической записи. Он не знает, что это ClickHouse.
	fmt.Printf("ИНФО: Передача %d валидных событий в хранилище...\n", len(validEvents))
	return s.writer.Batch(ctx, validEvents)
}
