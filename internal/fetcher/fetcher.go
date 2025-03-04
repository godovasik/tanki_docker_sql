package fetcher

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/godovasik/tanki_docker_sql/internal/models"
)

// Fetcher определяет интерфейс для запроса данных из API
type Fetcher interface {
	sendRequest(ctx context.Context, username string) (*http.Response, error)
	parseResponse(resp *http.Response) (*models.ResponseWrapper, error)
	GetUserStats(ctx context.Context, username string) (*models.ResponseWrapper, error)
}

// HTTPFetcher реализует Fetcher с использованием http.Client
type HTTPFetcher struct {
	client *http.Client
}

// NewHTTPFetcher создает новый HTTPFetcher с заданным таймаутом
func NewHTTPFetcher(timeout time.Duration) *HTTPFetcher {
	return &HTTPFetcher{
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Отключает проверку сертификата
			},
		},
	}
}

// SendRequest отправляет HTTP-запрос с учетом контекста
func (h *HTTPFetcher) sendRequest(ctx context.Context, username string) (*http.Response, error) {
	url := fmt.Sprintf("http://ratings.tankionline.com/api/eu/profile/?user=%s&lang=en", username)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// ParseResponse читает и парсит JSON-ответ
func (h *HTTPFetcher) parseResponse(resp *http.Response) (*models.ResponseWrapper, error) {
	var data models.ResponseWrapper
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &data, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return &data, err
	}

	return &data, nil
}

func (h *HTTPFetcher) GetUserStats(ctx context.Context, username string) (*models.ResponseWrapper, error) {
	resp, err := h.sendRequest(ctx, username) // надо обработать если юзер не найден
	if err != nil {
		return nil, fmt.Errorf("ошибка при отправке реквеста: %v", err)
	}
	// fmt.Println("resp:", resp)
	data, err := h.parseResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга: %v", err)
	}

	if data.ResponseType == "NOT_FOUND" {
		return nil, fmt.Errorf("NOT_FOUND")
	}
	return data, nil
}
