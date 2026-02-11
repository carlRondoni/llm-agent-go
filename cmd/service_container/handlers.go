package service_container

import "llm-agent-go/internal/application"

type Handlers struct {
	GenerateHandler  application.GenerateHandler
	StreamHandler    application.StreamHandler
	LLMHealthHandler application.LLMHealthCheckHandler
}

func NewHandlers(llmClients LLMClients) Handlers {
	return Handlers{
		GenerateHandler:  application.NewGenerateHandler(llmClients.OllamaClient),
		StreamHandler:    application.NewStreamHandler(llmClients.OllamaClient),
		LLMHealthHandler: application.NewLLMHealthCheckHandler(llmClients.OllamaClient),
	}
}
