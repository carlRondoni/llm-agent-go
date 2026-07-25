package controllers_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"llm-agent-go/internal/application"
	"llm-agent-go/internal/infrastructure/controllers"
)

type mockLLMClient struct {
	generateFn func(ctx context.Context, prompt string) (string, error)
	streamFn   func(ctx context.Context, prompt string) (<-chan string, error)
	healthFn   func(ctx context.Context) error
}

func (m mockLLMClient) Generate(ctx context.Context, prompt string) (string, error) {
	if m.generateFn != nil {
		return m.generateFn(ctx, prompt)
	}
	return "", nil
}

func (m mockLLMClient) Stream(ctx context.Context, prompt string) (<-chan string, error) {
	if m.streamFn != nil {
		return m.streamFn(ctx, prompt)
	}
	return nil, nil
}

func (m mockLLMClient) Health(ctx context.Context) error {
	if m.healthFn != nil {
		return m.healthFn(ctx)
	}
	return nil
}

func TestGenerateController_Success(t *testing.T) {
	mockClient := mockLLMClient{
		generateFn: func(ctx context.Context, prompt string) (string, error) {
			return "hello world", nil
		},
	}
	handler := application.NewGenerateHandler(mockClient)
	controller := controllers.NewGenerateController(handler, zerolog.Nop())

	body := []byte(`{"prompt": "hi"}`)
	req := httptest.NewRequest(http.MethodPost, "/llm/generate", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	controller.Execute(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "hello world")
}

func TestGenerateController_MethodNotAllowed(t *testing.T) {
	controller := controllers.NewGenerateController(application.GenerateHandler{}, zerolog.Nop())
	req := httptest.NewRequest(http.MethodGet, "/llm/generate", nil)
	w := httptest.NewRecorder()

	controller.Execute(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Result().StatusCode)
}

func TestGenerateController_InvalidBody(t *testing.T) {
	controller := controllers.NewGenerateController(application.GenerateHandler{}, zerolog.Nop())
	req := httptest.NewRequest(http.MethodPost, "/llm/generate", bytes.NewBuffer([]byte(`{bad-json}`)))
	w := httptest.NewRecorder()

	controller.Execute(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestGenerateController_EmptyPrompt(t *testing.T) {
	controller := controllers.NewGenerateController(application.GenerateHandler{}, zerolog.Nop())
	req := httptest.NewRequest(http.MethodPost, "/llm/generate", bytes.NewBuffer([]byte(`{"prompt": ""}`)))
	w := httptest.NewRecorder()

	controller.Execute(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestGenerateController_InternalError(t *testing.T) {
	mockClient := mockLLMClient{
		generateFn: func(ctx context.Context, prompt string) (string, error) {
			return "", errors.New("llm error")
		},
	}
	handler := application.NewGenerateHandler(mockClient)
	controller := controllers.NewGenerateController(handler, zerolog.Nop())

	body := []byte(`{"prompt": "hi"}`)
	req := httptest.NewRequest(http.MethodPost, "/llm/generate", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	controller.Execute(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}
