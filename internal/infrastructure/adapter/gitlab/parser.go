package gitlab

import (
	"fmt"
	"os"

	"github.com/nalonluap/ci-turbo/internal/app/service"
	"github.com/nalonluap/ci-turbo/internal/domain"
	"gopkg.in/yaml.v3"
)

// Убеждаемся, что Parser реализует интерфейс PipelineRepository.
var _ service.PipelineRepository = (*Parser)(nil)

// gitlabCIConfig — это структура, которая точно отражает формат .gitlab-ci.yml.
// Она используется как промежуточное представление перед трансформацией в доменную модель.
type gitlabCIConfig struct {
	Stages      []string             `yaml:"stages"`
	GlobalCache *gitlabCache         `yaml:"cache"`
	Jobs        map[string]gitlabJob `yaml:",inline"` // `inline` позволяет парсить все остальные ключи как jobs
}

// gitlabJob представляет одну задачу в .gitlab-ci.yml.
type gitlabJob struct {
	Stage    string       `yaml:"stage"`
	Script   interface{}  `yaml:"script"` // `script` может быть строкой или списком строк
	Cache    *gitlabCache `yaml:"cache"`
	Parallel int          `yaml:"parallel"`
}

// gitlabCache представляет блок `cache` в .gitlab-ci.yml.
type gitlabCache struct {
	Key   string   `yaml:"key"`
	Paths []string `yaml:"paths"`
}

// Parser - это адаптер, который читает .gitlab-ci.yml и превращает его
// в нашу универсальную доменную модель.
type Parser struct {
	// Здесь могут быть зависимости, например, логгер.
}

func NewParser() *Parser {
	return &Parser{}
}

// Get реализует логику получения и трансформации данных.
func (p *Parser) Get(path string) (*domain.Pipeline, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл %s: %w", path, err)
	}

	// 2. Парсим YAML в наши промежуточные gitlab-структуры
	var config gitlabCIConfig
	if err = yaml.Unmarshal(yamlFile, &config); err != nil {
		return nil, fmt.Errorf("не удалось распарсить YAML: %w", err)
	}

	// 3. Трансформируем промежуточные структуры в нашу универсальную доменную модель
	pipeline := mapToDomain(config)

	return pipeline, nil
}

// mapToDomain - это приватная функция-маппер.
// Она выполняет "грязную" работу по преобразованию одной структуры в другую.
func mapToDomain(config gitlabCIConfig) *domain.Pipeline {
	domainJobs := make(map[string]domain.Job)

	// Итерируемся по всем задачам, найденным в YAML
	for jobName, job := range config.Jobs {
		// Пропускаем "скрытые" задачи, которые начинаются с точки
		if jobName[0] == '.' {
			continue
		}

		// Создаем доменную модель задачи
		domainJob := domain.Job{
			Name:     jobName,
			Stage:    job.Stage,
			Script:   normalizeScript(job.Script),
			Parallel: job.Parallel,
		}

		// Определяем настройки кэша: приоритет у локального кэша задачи,
		// если его нет, используем глобальный.
		finalCache := job.Cache
		if finalCache == nil {
			finalCache = config.GlobalCache
		}

		if finalCache != nil {
			domainJob.Cache = domain.CacheConfig{
				Key:   finalCache.Key,
				Paths: finalCache.Paths,
			}
		}

		domainJobs[jobName] = domainJob
	}

	return &domain.Pipeline{
		Jobs: domainJobs,
	}
}

// normalizeScript — утилита для обработки поля 'script', которое может
// быть как одной строкой, так и списком строк.
func normalizeScript(script interface{}) []string {
	if script == nil {
		return nil
	}

	// Если это список строк, возвращаем его
	if lines, ok := script.([]interface{}); ok {
		result := make([]string, len(lines))
		for i, v := range lines {
			if str, ok := v.(string); ok {
				result[i] = str
			}
		}
		return result
	}

	// Если это одна строка, возвращаем срез из одного элемента
	if line, ok := script.(string); ok {
		return []string{line}
	}

	return nil
}
