package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/nalonluap/ci-turbo/internal/domain" // Импортируем общую доменную модель
	"github.com/nalonluap/ci-turbo/services/ingester/internal/interfaces"

	"github.com/nalonluap/ci-turbo/services/ingester/internal/config"

	"github.com/nalonluap/ci-turbo/services/ingester/internal/service"

	"github.com/nalonluap/ci-turbo/services/ingester/pkg/broker/nats"
	"github.com/nalonluap/ci-turbo/services/ingester/pkg/storage/clickhouse"
)

func main() {
	log.Println("Запуск Ingestion Worker...")

	// 1. Загружаем конфигурацию (адреса Kafka, DSN для ClickHouse и т.д.)
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// 2. Инициализируем адаптеры (инфраструктурный слой)

	// Консьюмер для чтения сообщений из Kafka
	consumer, err := nats.NewConsumer(cfg.NATS.URL, cfg.NATS.ConsumerGroup)
	if err != nil {
		log.Fatalf("Не удалось создать Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Репозиторий для записи данных в ClickHouse
	metricStorage, err := clickhouse.NewRepository(cfg.ClickHouse)
	if err != nil {
		log.Fatalf("Не удалось подключиться к ClickHouse: %v", err)
	}

	// 3. Инициализируем сервис (слой приложения), внедряя зависимость
	processingService := service.NewProcessingService(metricStorage)

	// 4. Запускаем бесконечный цикл обработки
	log.Printf("Воркер начал прослушивание топика '%s'...", cfg.NATS.StreamName)
	ctx := context.Background()

	for {
		// Читаем пачку сообщений из Kafka. Это блокирующая операция с таймаутом.
		messages, err := consumer.FetchBatch(ctx, 100, 5*time.Second)
		if err != nil {
			log.Printf("ОШИБКА: Не удалось получить пачку из Kafka: %v", err)
			time.Sleep(5 * time.Second) // Ждем немного перед следующей попыткой
			continue
		}

		if len(messages) == 0 {
			continue // Если сообщений нет, начинаем новый цикл ожидания
		}

		// Преобразуем сырые байты из Kafka в наши доменные структуры
		events := transformMessagesToDomain(messages)
		if len(events) == 0 {
			log.Println("ПРЕДУПРЕЖДЕНИЕ: В пачке не найдено валидных событий после парсинга.")
			// Подтверждаем сообщения, чтобы не обрабатывать "битые" данные вечно
			if err = consumer.CommitMessages(ctx, messages); err != nil {
				log.Printf("ОШИБКА: Не удалось подтвердить 'битые' сообщения в Kafka: %v", err)
			}
			continue
		}

		// Вызываем наш сервис для обработки и записи пачки
		if err = processingService.ProcessBatch(ctx, events); err != nil {
			log.Printf("ОШИБКА: Не удалось обработать пачку: %v. Сообщения будут обработаны повторно.", err)
			// Важно: мы НЕ подтверждаем (commit) сообщения, Kafka вернет их нам снова
			continue
		}

		// Если запись в БД прошла успешно, подтверждаем Kafka, что сообщения можно удалить из очереди
		if err = consumer.CommitMessages(ctx, messages); err != nil {
			log.Printf("ОШИБКА: Не удалось подтвердить сообщения в Kafka: %v", err)
		}

		log.Printf("Успешно обработано и сохранено %d событий.", len(events))
	}
}

// transformMessagesToDomain - это вспомогательная функция, которая преобразует
// сырые сообщения из брокера в типизированные структуры нашего домена.
func transformMessagesToDomain(messages []interfaces.Message) []*domain.Event {
	// Создаем срез с предвыделенной емкостью для эффективности
	events := make([]*domain.Event, 0, len(messages))

	for _, msg := range messages {
		var event domain.Event

		// Десериализуем JSON из тела сообщения, используя метод интерфейса
		if err := json.Unmarshal(msg.Payload(), &event); err != nil {
			log.Printf("ПРЕДУПРЕЖДЕНИЕ: Не удалось распарсить JSON из сообщения: %v", err)
			continue // Пропускаем "битое" сообщение
		}

		events = append(events, &event)
	}
	return events
}
