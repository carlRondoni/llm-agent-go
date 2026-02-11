package application

import (
	"context"
	"llm-agent-go/internal/domain"
)

type StreamHandler struct {
	llmClient domain.LLMClient
}

func NewStreamHandler(
	llmClient domain.LLMClient,
) StreamHandler {
	return StreamHandler{
		llmClient: llmClient,
	}
}

func (h StreamHandler) Handle(ctx context.Context) error {
	return nil
}
