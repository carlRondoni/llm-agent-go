package llm_clients

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sony/gobreaker/v2"
	"go.opentelemetry.io/otel"
	"golang.org/x/time/rate"
)

type OllamaClient struct {
	baseURL string
	model   string
	http    *http.Client
	breaker *gobreaker.CircuitBreaker[any]
	limiter *rate.Limiter
}

func NewOllamaClient(baseURL string, model string) OllamaClient {
	cbst := gobreaker.Settings{
		Name:        "LLM",
		MaxRequests: 1,
		Interval:    60 * time.Second,
		Timeout:     10 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 3
		},
	}

	cb := gobreaker.NewCircuitBreaker[any](cbst)

	return OllamaClient{
		baseURL: baseURL,
		model:   model,
		http: &http.Client{
			Timeout: 2 * time.Minute,
		},
		breaker: cb,
		limiter: rate.NewLimiter(2, 5),
	}
}

func (c OllamaClient) Generate(ctx context.Context, prompt string) (string, error) {
	log.Info().Str("prompt", prompt).Msg("generate")

	if err := c.limiter.Wait(ctx); err != nil {
		return "", err
	}

	tr := otel.Tracer("llm-api")
	ctx, span := tr.Start(ctx, "Generate")
	defer span.End()

	result, err := c.breaker.Execute(func() (any, error) {
		reqBody := map[string]string{
			"model":  c.model,
			"prompt": prompt,
		}

		b, _ := json.Marshal(reqBody)

		req, _ := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/generate", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.http.Do(req)
		if err != nil {
			log.Error().
				Err(err).
				Str("prompt", prompt).
				Msg("failed to call LLM")
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			return "", fmt.Errorf("llm generate failed, status: %d", resp.StatusCode)
		}

		var out struct {
			Response string `json:"response"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			log.Error().
				Err(err).
				Str("prompt", prompt).
				Msg("failed to call LLM")
			return nil, err
		}

		return out.Response, nil
	})

	if err != nil {
		log.Error().
			Err(err).
			Str("prompt", prompt).
			Msg("failed to call LLM")
		return "", err
	}

	response, ok := result.(string)
	if !ok {
		return "", fmt.Errorf("unexpected response type: %T", result)
	}

	return response, nil
}

func (c OllamaClient) Stream(ctx context.Context, prompt string) (<-chan string, error) {
	log.Info().Str("prompt", prompt).Msg("stream")

	if err := c.limiter.Wait(ctx); err != nil {
		return nil, err
	}

	tr := otel.Tracer("llm-api")
	ctx, span := tr.Start(ctx, "Generate")
	defer span.End()

	out := make(chan string)

	_, err := c.breaker.Execute(func() (any, error) {
		body := map[string]any{
			"model":  c.model,
			"prompt": prompt,
			"stream": true,
		}

		b, _ := json.Marshal(body)

		req, _ := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/generate", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.http.Do(req)
		if err != nil {
			log.Error().
				Err(err).
				Str("prompt", prompt).
				Msg("failed to call LLM")
			return nil, err
		}

		out := make(chan string)

		go func() {
			defer resp.Body.Close()
			defer close(out)

			scanner := bufio.NewScanner(resp.Body)

			for scanner.Scan() {
				var chunk struct {
					Response string `json:"response"`
					Done     bool   `json:"done"`
				}

				if err := json.Unmarshal(scanner.Bytes(), &chunk); err != nil {
					log.Error().
						Err(err).
						Str("prompt", prompt).
						Msg("failed to call LLM")
					continue
				}

				if chunk.Response != "" {
					out <- chunk.Response
				}

				if chunk.Done {
					return
				}
			}
		}()

		return nil, nil
	})
	if err != nil {
		log.Error().
			Err(err).
			Str("prompt", prompt).
			Msg("failed to call LLM")
		close(out)
		return nil, err
	}

	return out, nil
}

func (c OllamaClient) Health(ctx context.Context) error {
	log.Info().Msg("health")

	if err := c.limiter.Wait(ctx); err != nil {
		return err
	}

	tr := otel.Tracer("llm-api")
	ctx, span := tr.Start(ctx, "Generate")
	defer span.End()

	_, err := c.breaker.Execute(func() (any, error) {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			c.baseURL+"/api/tags",
			nil,
		)
		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to call LLM")
			return nil, err
		}

		resp, err := c.http.Do(req)
		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to call LLM")
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 500 {
			return nil, fmt.Errorf("llm unhealthy status: %d", resp.StatusCode)
		}

		return nil, nil
	})

	return err
}
