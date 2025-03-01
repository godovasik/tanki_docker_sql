package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/godovasik/tanki_docker_sql/internal/models"
)

// Fetcher определяет интерфейс для запроса данных из API
type Fetcher interface {
	SendRequest(ctx context.Context, username string) (*http.Response, error)
	ParseResponse(resp *http.Response) (*models.ResponseWrapper, error)
}

// HTTPFetcher реализует Fetcher с использованием http.Client
type HTTPFetcher struct {
	client *http.Client
}

// NewHTTPFetcher создает новый HTTPFetcher с заданным таймаутом
func NewHTTPFetcher(timeout time.Duration) *HTTPFetcher {
	return &HTTPFetcher{
		client: &http.Client{Timeout: timeout},
	}
}

// SendRequest отправляет HTTP-запрос с учетом контекста
func (h *HTTPFetcher) SendRequest(ctx context.Context, username string) (*http.Response, error) {
	url := fmt.Sprintf("https://ratings.tankionline.com/api/eu/profile/?user=%s&lang=en", username)

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
func (h *HTTPFetcher) ParseResponse(resp *http.Response) (*models.ResponseWrapper, error) {
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
