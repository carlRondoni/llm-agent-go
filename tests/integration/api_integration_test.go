//go:build integration

package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getAPIBaseURL(t *testing.T) string {
	t.Helper()

	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	client := &http.Client{Timeout: 1 * time.Second}
	require.Eventually(t, func() bool {
		resp, err := client.Get(apiURL + "/health")
		if err != nil {
			return false
		}
		defer func() { assert.NoError(t, resp.Body.Close()) }()
		return resp.StatusCode == http.StatusOK
	}, 30*time.Second, 500*time.Millisecond, "API must be running at %s (start via make docker-start-build)", apiURL)

	return apiURL
}

func postPrompt(t *testing.T, client *http.Client, url string) *http.Response {
	t.Helper()

	b, err := json.Marshal(map[string]string{"prompt": "Reply with one short sentence about Go."})
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewBuffer(b))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	require.NoError(t, err)

	return resp
}

func TestAPI_E2E_Integration(t *testing.T) {
	baseURL := getAPIBaseURL(t)
	client := &http.Client{Timeout: 90 * time.Second}

	t.Run("GET /health", func(t *testing.T) {
		resp, err := client.Get(baseURL + "/health")
		require.NoError(t, err)
		defer func() { assert.NoError(t, resp.Body.Close()) }()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET /llm/health", func(t *testing.T) {
		resp, err := client.Get(baseURL + "/llm/health")
		require.NoError(t, err)
		defer func() { assert.NoError(t, resp.Body.Close()) }()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body map[string]string
		err = json.NewDecoder(resp.Body).Decode(&body)
		require.NoError(t, err)
		assert.Equal(t, "ok", body["status"])
		assert.Equal(t, "up", body["llm"])
	})

	// small models don't follow instructions reliably, so only check that text comes back
	t.Run("POST /llm/generate", func(t *testing.T) {
		resp := postPrompt(t, client, baseURL+"/llm/generate")
		defer func() { assert.NoError(t, resp.Body.Close()) }()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result struct {
			Response string `json:"response"`
		}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.NotEmpty(t, strings.TrimSpace(result.Response))
	})

	t.Run("POST /llm/stream", func(t *testing.T) {
		resp := postPrompt(t, client, baseURL+"/llm/stream")
		defer func() { assert.NoError(t, resp.Body.Close()) }()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, []string{"chunked"}, resp.TransferEncoding)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.NotEmpty(t, strings.TrimSpace(string(body)))
	})
}

// validation runs before any call to Ollama, so these answers are deterministic
func TestAPI_Errors_Integration(t *testing.T) {
	baseURL := getAPIBaseURL(t)
	client := &http.Client{Timeout: 10 * time.Second}

	tests := []struct {
		name   string
		method string
		path   string
		body   string
		status int
		text   string
	}{
		{"generate with invalid json", http.MethodPost, "/llm/generate", `{bad`, http.StatusBadRequest, "invalid json body"},
		{"generate with empty prompt", http.MethodPost, "/llm/generate", `{"prompt":""}`, http.StatusBadRequest, "prompt required"},
		{"stream with invalid json", http.MethodPost, "/llm/stream", `{bad`, http.StatusBadRequest, "invalid json body"},
		{"stream without prompt", http.MethodPost, "/llm/stream", `{}`, http.StatusBadRequest, "prompt required"},
		{"GET /llm/generate", http.MethodGet, "/llm/generate", "", http.StatusMethodNotAllowed, "method not allowed"},
		{"GET /llm/stream", http.MethodGet, "/llm/stream", "", http.StatusMethodNotAllowed, "method not allowed"},
		{"POST /llm/health", http.MethodPost, "/llm/health", "", http.StatusMethodNotAllowed, "method not allowed"},
		{"POST /health", http.MethodPost, "/health", "", http.StatusMethodNotAllowed, "method not allowed"},
		{"old /health/llm route", http.MethodGet, "/health/llm", "", http.StatusNotFound, "not found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequestWithContext(context.Background(), tt.method, baseURL+tt.path, strings.NewReader(tt.body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer func() { assert.NoError(t, resp.Body.Close()) }()

			assert.Equal(t, tt.status, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Contains(t, string(body), tt.text)
		})
	}
}
