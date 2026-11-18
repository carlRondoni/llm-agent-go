//go:build integration

package llm_clients_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"llm-agent-go/internal/infrastructure/llm_clients"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestOllamaClient_Integration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	baseURL := getOllamaURL()
	model := getOllamaModel()
	client := llm_clients.NewOllamaClient(baseURL, model, zerolog.Nop())

	err := client.Health(ctx)
	require.NoError(t, err, "Ollama must be running at %s for integration tests (start via make docker-start)", baseURL)

	// small models don't follow instructions reliably, so only check that text comes back
	response, err := client.Generate(ctx, "Reply with one short sentence about Go.")
	require.NoError(t, err, "failed to generate text with model %s", model)
	assert.NotEmpty(t, strings.TrimSpace(response))

	stream, err := client.Stream(ctx, "Reply with one short sentence about Go.")
	require.NoError(t, err, "failed to start a stream with model %s", model)

	var chunks int
	var text strings.Builder
	for chunk := range stream {
		chunks++
		text.WriteString(chunk)
	}
	assert.Positive(t, chunks)
	assert.NotEmpty(t, strings.TrimSpace(text.String()))
}

// Ollama answers 404 for a model it doesn't have
func TestOllamaClient_UnknownModel_Integration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := llm_clients.NewOllamaClient(getOllamaURL(), "llm-agent-go-missing-model", zerolog.Nop())

	_, err := client.Generate(ctx, "hi")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "status: 404")

	// documented limitation: Stream doesn't check the status, so the channel just closes empty
	stream, err := client.Stream(ctx, "hi")
	require.NoError(t, err)

	var chunks int
	for range stream {
		chunks++
	}
	assert.Zero(t, chunks)
}
