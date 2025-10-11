package api

import (
	"fmt"

	"github.com/nalonluap/ci-turbo/services/observer/internal/service"

	"github.com/gin-gonic/gin"
)

// Server — это обертка для нашего HTTP-сервера, хранящая зависимости.
type Server struct {
	port          string
	engine        *gin.Engine
	ingestService *service.IngestService
	queryService  *service.QueryService // Для будущих аналитических эндпоинтов
}

// NewServer создает и настраивает новый экземпляр сервера.
func NewServer(port string, ingestService *service.IngestService, queryService *service.QueryService) *Server {
	server := &Server{
		port:          port,
		engine:        gin.Default(), // Создаем роутер Gin со стандартными middleware (логгер, recovery)
		ingestService: ingestService,
		queryService:  queryService,
	}

	server.Routes() // Вызываем настройку маршрутов

	return server
}

// Routes определяет все маршруты (эндпоинты) для нашего API.
func (s *Server) Routes() {
	v1 := s.engine.Group("/v1")
	{
		// 1. Хендлер для ПРИЕМА данных от агентов
		v1.POST("/metrics", s.handleIngestMetric)

		// 2. Хендлер для ПРЕДОСТАВЛЕНИЯ аналитики дашборду
		v1.GET("/analytics/slowest-jobs", s.handleGetSlowestJobs)
	}
}

func (s *Server) Run() error {
	return s.engine.Run(fmt.Sprintf(":%s", s.port))
}
