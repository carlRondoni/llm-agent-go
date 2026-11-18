//go:build integration

package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"llm-agent-go/internal/application"
	"llm-agent-go/internal/infrastructure/controllers"
	"llm-agent-go/internal/infrastructure/llm_clients"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// These tests build the API in the test process, with the real controllers, handlers and
// Ollama client against the real Ollama. That makes it possible to pick the model and to
// watch the requests that go to Ollama.

func getOllamaURL() string {
	if url := os.Getenv("OLLAMA_URL"); url != "" {
		return url
	}
	return "http://localhost:11434"
}

func getOllamaModel() string {
	if model := os.Getenv("OLLAMA_MODEL"); model != "" {
		return model
	}
	return "llama3"
}

func newInProcessAPI(t *testing.T, ollamaURL string, model string) *httptest.Server {
	t.Helper()

	client := llm_clients.NewOllamaClient(ollamaURL, model, zerolog.Nop())

	mux := http.NewServeMux()
	mux.HandleFunc("/llm/generate", controllers.NewGenerateController(application.NewGenerateHandler(client), zerolog.Nop()).Execute)
	mux.HandleFunc("/llm/stream", controllers.NewStreamController(application.NewStreamHandler(client), zerolog.Nop()).Execute)

	api := httptest.NewServer(mux)
	t.Cleanup(api.Close)

	return api
}

func newPromptRequest(t *testing.T, ctx context.Context, url string, prompt string) *http.Request {
	t.Helper()

	b, err := json.Marshal(map[string]string{"prompt": prompt})
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(b))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	return req
}

// Ollama answers 404 for a model it doesn't have
func TestAPI_UnknownModel_Integration(t *testing.T) {
	api := newInProcessAPI(t, getOllamaURL(), "llm-agent-go-missing-model")
	client := &http.Client{Timeout: 30 * time.Second}

	t.Run("generate answers 500 with the error", func(t *testing.T) {
		resp, err := client.Do(newPromptRequest(t, context.Background(), api.URL+"/llm/generate", "hi"))
		require.NoError(t, err)
		defer func() { assert.NoError(t, resp.Body.Close()) }()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var body map[string]string
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.Contains(t, body["error"], "status: 404")
	})

	// documented limitation: the stream doesn't check Ollama's status code
	t.Run("stream answers an empty 200", func(t *testing.T) {
		resp, err := client.Do(newPromptRequest(t, context.Background(), api.URL+"/llm/stream", "hi"))
		require.NoError(t, err)
		defer func() { assert.NoError(t, resp.Body.Close()) }()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Empty(t, string(body))
	})
}

// newObservingProxy forwards every request to Ollama unchanged, and reports for each
// /api/generate request whether Ollama's response was read to the end or cut off.
func newObservingProxy(t *testing.T, ollamaURL string) (*httptest.Server, <-chan bool) {
	t.Helper()

	completed := make(chan bool, 10)

	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequestWithContext(r.Context(), r.Method, ollamaURL+r.URL.Path, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		req.Header = r.Header.Clone()

		resp, err := http.DefaultTransport.RoundTrip(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer func() { _ = resp.Body.Close() }()

		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(resp.StatusCode)

		buf := make([]byte, 4096)
		for {
			n, readErr := resp.Body.Read(buf)
			if n > 0 {
				if _, err := w.Write(buf[:n]); err != nil {
					break
				}
				w.(http.Flusher).Flush()
			}
			if readErr != nil {
				if r.URL.Path == "/api/generate" {
					completed <- readErr == io.EOF
				}
				return
			}
		}
		if r.URL.Path == "/api/generate" {
			completed <- false
		}
	}))
	t.Cleanup(proxy.Close)

	return proxy, completed
}

func TestAPI_StreamDisconnect_CancelsOllamaRequest_Integration(t *testing.T) {
	proxy, completed := newObservingProxy(t, getOllamaURL())
	api := newInProcessAPI(t, proxy.URL, getOllamaModel())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// long enough that the model is still generating when the client leaves
	prompt := "Write a very long story, at least 1000 words, about a lighthouse keeper."
	resp, err := http.DefaultClient.Do(newPromptRequest(t, ctx, api.URL+"/llm/stream", prompt))
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	first := make([]byte, 1)
	_, err = io.ReadFull(resp.Body, first)
	require.NoError(t, err, "expected at least one chunk before disconnecting")

	// the client goes away in the middle of the stream
	cancel()

	select {
	case done := <-completed:
		assert.False(t, done, "Ollama finished the answer, so the test couldn't observe a cancellation; use a longer prompt")
	case <-time.After(10 * time.Second):
		t.Fatal("the request to Ollama was still open 10 seconds after the client disconnected")
	}

	// the service keeps working after the disconnect
	client := &http.Client{Timeout: 90 * time.Second}
	next, err := client.Do(newPromptRequest(t, context.Background(), api.URL+"/llm/generate", "Reply with one short sentence about Go."))
	require.NoError(t, err)
	defer func() { assert.NoError(t, next.Body.Close()) }()
	assert.Equal(t, http.StatusOK, next.StatusCode)

	var result struct {
		Response string `json:"response"`
	}
	require.NoError(t, json.NewDecoder(next.Body).Decode(&result))
	assert.NotEmpty(t, strings.TrimSpace(result.Response))
}
