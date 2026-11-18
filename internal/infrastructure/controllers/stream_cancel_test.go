//go:build unit || !integration

package controllers_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"llm-agent-go/internal/application"
	"llm-agent-go/internal/infrastructure/controllers"
	"llm-agent-go/internal/infrastructure/llm_clients"
)

// The whole chain is real here: HTTP server -> controller -> handler -> OllamaClient -> fake Ollama.
// Only Ollama is faked, so the test proves that a client disconnect reaches the upstream request.
func TestStreamController_ClientDisconnectCancelsUpstream(t *testing.T) {
	upstreamCancelled := make(chan struct{})

	ollama := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-ndjson")
		_, err := io.WriteString(w, `{"response":"hel","done":false}`+"\n")
		assert.NoError(t, err)
		w.(http.Flusher).Flush()

		// a model that keeps generating: only a cancelled request ends this
		select {
		case <-r.Context().Done():
			close(upstreamCancelled)
		case <-time.After(10 * time.Second):
		}
	}))
	defer ollama.Close()

	client := llm_clients.NewOllamaClient(ollama.URL, "test-model", zerolog.Nop())
	controller := controllers.NewStreamController(application.NewStreamHandler(client), zerolog.Nop())

	api := httptest.NewServer(http.HandlerFunc(controller.Execute))
	defer api.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, api.URL, bytes.NewBufferString(`{"prompt":"hi"}`))
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	buf := make([]byte, 3)
	_, err = io.ReadFull(resp.Body, buf)
	require.NoError(t, err)
	assert.Equal(t, "hel", string(buf))

	// the client goes away in the middle of the stream
	cancel()

	select {
	case <-upstreamCancelled:
	case <-time.After(2 * time.Second):
		t.Fatal("the request to Ollama was not cancelled after the client disconnected")
	}
}
