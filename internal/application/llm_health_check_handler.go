package application

import (
	"context"
	"llm-agent-go/internal/domain"
)

type LLMHealthCheckHandler struct {
	llmClient domain.LLMClient
}

func NewLLMHealthCheckHandler(
	llmClient domain.LLMClient,
) LLMHealthCheckHandler {
	return LLMHealthCheckHandler{
		llmClient: llmClient,
	}
}

func (h LLMHealthCheckHandler) Handle(ctx context.Context) error {
	return nil
}
