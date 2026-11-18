package service_container

import (
	"llm-agent-go/internal/infrastructure/llm_clients"
	"os"

	"github.com/rs/zerolog"
)

const defaultOllamaModel = "llama3"

type LLMClients struct {
	OllamaClient llm_clients.OllamaClient
}

func NewLLMClients(logger zerolog.Logger) LLMClients {
	url := os.Getenv("OLLAMA_URL")
	if url == "" {
		logger.Fatal().Msg("OLLAMA_URL is not set")
	}

	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		model = defaultOllamaModel
	}

	logger.Info().Str("model", model).Msg("using ollama model")

	return LLMClients{
		OllamaClient: llm_clients.NewOllamaClient(url, model, logger),
	}
}
