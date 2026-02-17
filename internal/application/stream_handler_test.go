package application_test

import (
	"context"
	"errors"
	"testing"

	"llm-agent-go/internal/application"

	"github.com/stretchr/testify/assert"
)

type mockLLMStream struct {
	called bool
	ctx    context.Context
	prompt string
	ch     <-chan string
	err    error
}

func (m *mockLLMStream) Stream(ctx context.Context, prompt string) (<-chan string, error) {
	m.called = true
	m.ctx = ctx
	m.prompt = prompt
	return m.ch, m.err
}

func (m *mockLLMStream) Generate(context.Context, string) (string, error) {
	panic("not used")
}

func (m *mockLLMStream) Health(context.Context) error {
	panic("not used")
}

func TestStreamHandler_CallsClient(t *testing.T) {
	out := make(chan string)
	mock := &mockLLMStream{ch: out}

	h := application.NewStreamHandler(mock)

	ctx := context.WithValue(context.Background(), "k", "v")
	res, err := h.Handle(ctx, "hello")
	assert.NoError(t, err)
	assert.True(t, mock.called)
	assert.Equal(t, "hello", mock.prompt)
	assert.Equal(t, "v", mock.ctx.Value("k"))
	assert.Equal(t, res, (<-chan string)(out))
}

func TestStreamHandler_ReturnsError(t *testing.T) {
	expected := errors.New("boom")

	mock := &mockLLMStream{
		err: expected,
	}

	h := application.NewStreamHandler(mock)

	_, err := h.Handle(context.Background(), "x")
	assert.Error(t, err)
	assert.Equal(t, expected, err)
}
