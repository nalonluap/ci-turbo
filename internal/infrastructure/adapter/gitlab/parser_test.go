package gitlab_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nalonluap/ci-turbo/internal/app/service"
	"github.com/nalonluap/ci-turbo/internal/infrastructure/adapter/gitlab"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Эта статическая проверка, которую мы обсуждали, остается здесь.
// Она гарантирует, что наш Parser не "сломал" контракт с интерфейсом.
func TestParserImplementsPipelineRepository(t *testing.T) {
	var _ service.PipelineRepository = (*gitlab.Parser)(nil)
}

// TestParser_Get - это основной набор тестов для метода Get.
func TestParser_Get(t *testing.T) {
	// Определяем структуру для наших тестовых случаев
	testCases := []struct {
		name              string // Название теста
		fileContent       string // Содержимое .gitlab-ci.yml, которое мы будем создавать
		expectErr         bool   // Ожидаем ли мы ошибку
		expectedJobsCount int    // Ожидаемое количество задач в пайплайне
		expectedJobName   string // Имя одной из задач для выборочной проверки
	}{
		{
			name: "Valid and simple config",
			fileContent: `
stages:
  - build
  - test

build_job:
  stage: build
  script:
    - echo "Building..."

test_job:
  stage: test
  script:
    - echo "Testing..."
`,
			expectErr:         false,
			expectedJobsCount: 2,
			expectedJobName:   "build_job",
		},
		{
			name: "Config with cache definition",
			fileContent: `
build_app:
  script:
    - npm install
  cache:
    key: "node_modules"
    paths:
      - node_modules/
`,
			expectErr:         false,
			expectedJobsCount: 1,
			expectedJobName:   "build_app",
		},
		{
			name:              "Empty file should not cause panic",
			fileContent:       ``,
			expectErr:         false, // Пустой файл - это валидный YAML
			expectedJobsCount: 0,
		},
		{
			name:        "Invalid YAML syntax",
			fileContent: `job: [invalid syntax`,
			expectErr:   true,
		},
		{
			name:        "Non-existent file",
			fileContent: ``, // Контент не важен, так как файла не будет
			expectErr:   true,
		},
	}

	// Запускаем цикл по всем нашим тестовым случаям
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := gitlab.NewParser()
			filePath := ""

			if tc.name != "Non-existent file" {
				// Создаем временную директорию и файл для теста
				tempDir, err := os.MkdirTemp("", "ci-turbo-tests")
				require.NoError(t, err)
				defer os.RemoveAll(tempDir) // Гарантируем удаление после теста

				filePath = filepath.Join(tempDir, ".gitlab-ci.yml")
				err = os.WriteFile(filePath, []byte(tc.fileContent), 0644)
				require.NoError(t, err)
			} else {
				filePath = "non_existent_file.yml"
			}

			// Выполняем тестируемый метод
			pipeline, err := parser.Get(filePath)

			// Проверяем результат
			if tc.expectErr {
				// Если мы ожидали ошибку, убеждаемся, что она есть
				require.Error(t, err)
			} else {
				// Если ошибки не ожидали, убеждаемся, что ее нет
				require.NoError(t, err)
				// Убеждаемся, что сам пайплайн не nil
				require.NotNil(t, pipeline)
				// Проверяем, что количество задач совпадает с ожидаемым
				assert.Len(t, pipeline.Jobs, tc.expectedJobsCount)

				if tc.expectedJobName != "" {
					// Если нужно, делаем более глубокую проверку
					_, exists := pipeline.Jobs[tc.expectedJobName]
					assert.True(t, exists, "Expected job '%s' not found in pipeline", tc.expectedJobName)
				}
			}
		})
	}
}
