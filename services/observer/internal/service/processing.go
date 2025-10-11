package service

import (
	"context"
	"fmt"

	"github.com/nalonluap/ci-turbo/internal/domain"
	"github.com/nalonluap/ci-turbo/services/observer/internal/interfaces" // Воркер будет иметь свои интерфейсы
)

// ProcessingService отвечает за обработку пачек метрик и их сохранение в хранилище.
type ProcessingService struct {
	writer interfaces.Writer
}

// NewProcessingService создает новый экземпляр ProcessingService.
func NewProcessingService(writer interfaces.Writer) *ProcessingService {
	return &ProcessingService{writer: writer}
}

// Batch - это основной метод, который вызывает воркер.
// Он содержит логику фильтрации, обогащения и записи данных.
func (s *ProcessingService) Batch(ctx context.Context, events []*domain.Event) error {
	// 1. Фильтруем пустые или невалидные события.
	// Это пример бизнес-логики, которая живет именно в этом сервисе.
	validEvents := make([]*domain.Event, 0, len(events))
	for _, event := range events {
		if event.ProjectID == "" || event.JobID == "" {
			fmt.Printf("Пропущено невалидное событие: %+v\n", event)
			continue
		}
		// Здесь можно добавить логику обогащения данных
		validEvents = append(validEvents, event)
	}

	if len(validEvents) == 0 {
		return nil // Нечего записывать
	}

	// 2. Делегируем запись нашему хранилищу через интерфейс.
	return s.writer.Batch(ctx, validEvents)
}
