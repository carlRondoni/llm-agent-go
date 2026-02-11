package domain

import "context"

type LLMClient interface {
	Generate(ctx context.Context, prompt string) (string, error)
	Stream(ctx context.Context, prompt string) (<-chan string, error)
	Health(ctx context.Context) error
}
