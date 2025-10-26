package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nalonluap/ci-turbo/internal/domain"
)

func (s *Server) handleIngestMetric(c *gin.Context) {
	var event domain.Event

	// 1. Парсим и валидируем JSON
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// TODO: Добавить проверку API-токена из заголовков

	// 2. Вызываем сервис для обработки
	if err := s.ingestService.Ingest(c.Request.Context(), &event); err != nil {
		// Если брокер недоступен, возвращаем ошибку
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Failed to process event"})
		return
	}

	// 3. Отвечаем немедленно
	c.JSON(http.StatusAccepted, gin.H{"status": "event accepted"})
}

func (s *Server) handleGetSlowestJobs(c *gin.Context) {
	// TODO: Получить projectID из токена аутентификации или параметра запроса
	projectID := c.Query("projectID")

	// Важно: Добавляем проверку, что projectID был передан
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'projectID' is required"})
		return
	}

	// Здесь можно получить период из query-параметров, но пока используем значение по умолчанию
	results, err := s.queryService.SlowestJobs(c.Request.Context(), projectID, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch analytics"})
		return
	}

	c.JSON(http.StatusOK, results)
}
