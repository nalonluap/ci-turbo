package service

import (
	"context"
	"time"

	"github.com/nalonluap/ci-turbo/internal/domain"
	"github.com/nalonluap/ci-turbo/services/observer/internal/interfaces"
)

// QueryService инкапсулирует логику для выполнения аналитических запросов.
type QueryService struct {
	reader interfaces.Reader
}

// NewQueryService создает новый экземпляр QueryService.
func NewQueryService(reader interfaces.Reader) *QueryService {
	return &QueryService{reader: reader}
}

// SlowestJobs - это основной метод, который вызывает дашборд.
// Он содержит бизнес-логику, например, определение периода по умолчанию.
func (s *QueryService) SlowestJobs(ctx context.Context, projectID string, period time.Duration) ([]domain.SlowestJobAnalytics, error) {
	// Устанавливаем период по умолчанию, например, 7 дней.
	// Эта логика находится в сервисе, а не в хендлере или репозитории.
	defaultPeriod := 7 * 24 * time.Hour

	// Делегируем выполнение запроса нашему хранилищу через интерфейс.
	return s.reader.SlowestJobs(ctx, projectID, defaultPeriod)
}
