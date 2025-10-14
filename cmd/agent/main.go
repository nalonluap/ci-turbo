package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/nalonluap/ci-turbo/internal/domain"
)

func main() {
	if len(os.Args) < 3 || os.Args[1] != "measure" || os.Args[2] != "--" {
		log.Fatalf("Использование: ci-turbo-agent measure -- <команда> [аргументы...]")
	}

	// 1. Собираем команду и аргументы
	command := os.Args[3]
	args := os.Args[4:]

	// 2. Сбор метаданных из переменных окружения
	event := domain.Event{
		ProjectID:  os.Getenv("CI_PROJECT_ID"),
		PipelineID: os.Getenv("CI_PIPELINE_ID"),
		JobID:      os.Getenv("CI_JOB_ID"),
		JobName:    os.Getenv("CI_JOB_NAME"),
		RunnerID:   os.Getenv("CI_RUNNER_ID"),
		StartTime:  time.Now(),
	}

	// Вычисляем время в очереди
	if jobStartedAtStr := os.Getenv("CI_JOB_STARTED_AT"); jobStartedAtStr != "" {
		if jobStartedAt, err := time.Parse(time.RFC3339, jobStartedAtStr); err == nil {
			event.QueueTimeMs = time.Since(jobStartedAt).Milliseconds()
		}
	}

	// 3. Запуск дочернего процесса
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	startTime := time.Now()
	err := cmd.Run() // Блокируемся и ждем завершения
	duration := time.Since(startTime)

	// 4. Сбор результатов
	event.DurationMs = duration.Milliseconds()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			event.ExitCode = exitError.ExitCode()
		}
	} else {
		event.ExitCode = 0
	}

	// 5. Асинхронная отправка данных
	go sendMetrics(event)

	// 6. Возвращаем оригинальный код завершения
	os.Exit(event.ExitCode)
}

func sendMetrics(event domain.Event) {
	endpoint := os.Getenv("TURBO_OBSERVER_ENDPOINT")
	token := os.Getenv("TURBO_API_TOKEN")
	if endpoint == "" || token == "" {
		log.Println("ПРЕДУПРЕЖДЕНИЕ: Переменные TURBO_OBSERVER_ENDPOINT или TURBO_API_TOKEN не установлены. Метрики не будут отправлены.")
		return
	}

	payload, err := json.Marshal(event)
	if err != nil {
		log.Printf("ОШИБКА: Не удалось сериализовать метрики: %v", err)
		return
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("ОШИБКА: Не удалось создать запрос: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Turbo-Token", token)

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("ОШИБКА: Не удалось отправить метрики: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		log.Printf("ОШИБКА: Сервер ответил со статусом %d", resp.StatusCode)
	}
}
