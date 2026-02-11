package application

import (
	"context"
	"llm-agent-go/internal/domain"
)

type QueryHandler struct {
	llmClient domain.LLMClient
}

func NewQueryHandler(
	llmClient domain.LLMClient,
) QueryHandler {
	return QueryHandler{
		llmClient: llmClient,
	}
}

func (h QueryHandler) Handle(ctx context.Context) error {
	return nil
}
