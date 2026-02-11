package service_container

import (
	"llm-agent-go/internal/infrastructure/llm_clients"
	"log"
	"os"
)

type LLMClients struct {
	OllamaClient llm_clients.OllamaClient
}

func NewLLMClients() LLMClients {
	url := os.Getenv("OLLAMA_URL")
	if url == "" {
		log.Fatal("OLLAMA_URL is empty")
	}

	return LLMClients{
		OllamaClient: llm_clients.NewOllamaClient(url, "llama3"),
	}
}
