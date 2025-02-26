package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/godovasik/tanki_docker_sql/internal/models"
)

func SendRequest(username string) (*http.Response, error) {
	url := fmt.Sprintf("https://ratings.tankionline.com/api/eu/profile/?user=%s&lang=en", username)

	// Send GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}

	return resp, err
}

func ParseResponse(resp *http.Response) (models.ResponseWrapper, error) {
	var data models.ResponseWrapper
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}

	return data, err
}
