package service_container

import "llm-agent-go/internal/application"

type Handlers struct {
	QueryHandler     application.QueryHandler
	StreamHandler    application.StreamHandler
	LLMHealthHandler application.LLMHealthCheckHandler
}

func NewHandlers(llmClients LLMClients) Handlers {
	return Handlers{
		QueryHandler:     application.NewQueryHandler(llmClients.OllamaClient),
		StreamHandler:    application.NewStreamHandler(llmClients.OllamaClient),
		LLMHealthHandler: application.NewLLMHealthCheckHandler(llmClients.OllamaClient),
	}
}
