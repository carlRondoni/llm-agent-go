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

func TestStreamController_Success(t *testing.T) {
	mockClient := mockLLMClient{
		streamFn: func(ctx context.Context, prompt string) (<-chan string, error) {
			ch := make(chan string, 1)
			ch <- "chunk1"
			close(ch)
			return ch, nil
		},
	}
	handler := application.NewStreamHandler(mockClient)
	controller := controllers.NewStreamController(handler, zerolog.Nop())

	body := []byte(`{"prompt": "hi"}`)
	req := httptest.NewRequest(http.MethodPost, "/llm/stream", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	controller.Execute(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "chunk1")
}

func TestStreamController_MethodNotAllowed(t *testing.T) {
	controller := controllers.NewStreamController(application.StreamHandler{}, zerolog.Nop())
	req := httptest.NewRequest(http.MethodGet, "/llm/stream", nil)
	w := httptest.NewRecorder()

	controller.Execute(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Result().StatusCode)
}

func TestStreamController_InvalidBody(t *testing.T) {
	controller := controllers.NewStreamController(application.StreamHandler{}, zerolog.Nop())
	req := httptest.NewRequest(http.MethodPost, "/llm/stream", bytes.NewBuffer([]byte(`{bad-json}`)))
	w := httptest.NewRecorder()

	controller.Execute(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestStreamController_EmptyPrompt(t *testing.T) {
	controller := controllers.NewStreamController(application.StreamHandler{}, zerolog.Nop())
	req := httptest.NewRequest(http.MethodPost, "/llm/stream", bytes.NewBuffer([]byte(`{"prompt": ""}`)))
	w := httptest.NewRecorder()

	controller.Execute(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestStreamController_InternalError(t *testing.T) {
	mockClient := mockLLMClient{
		streamFn: func(ctx context.Context, prompt string) (<-chan string, error) {
			return nil, errors.New("stream error")
		},
	}
	handler := application.NewStreamHandler(mockClient)
	controller := controllers.NewStreamController(handler, zerolog.Nop())

	body := []byte(`{"prompt": "hi"}`)
	req := httptest.NewRequest(http.MethodPost, "/llm/stream", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	controller.Execute(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}
