package llm_clients

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type OllamaClient struct {
	baseURL string
	http    *http.Client
}

func NewOllamaClient(baseURL string) OllamaClient {
	return OllamaClient{
		baseURL: baseURL,
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c OllamaClient) Generate(ctx context.Context, prompt string) (string, error) {

	reqBody := map[string]string{
		"model":  "llama3",
		"prompt": prompt,
	}

	b, _ := json.Marshal(reqBody)

	req, _ := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/generate", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var out struct {
		Response string `json:"response"`
	}
	json.NewDecoder(resp.Body).Decode(&out)

	return out.Response, nil
}

func (c OllamaClient) Stream(ctx context.Context, prompt string) (<-chan string, error) {
	return nil, nil
}
