//go:build unit || !integration

package controllers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"context"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"llm-agent-go/internal/application"
	"llm-agent-go/internal/infrastructure/controllers"
)

func TestLLMHealthCheckController_Success(t *testing.T) {
	mockClient := mockLLMClient{
		healthFn: func(ctx context.Context) error {
			return nil
		},
	}
	handler := application.NewLLMHealthCheckHandler(mockClient)
	controller := controllers.NewLlmHealthCheckController(handler, zerolog.Nop())

	req := httptest.NewRequest(http.MethodGet, "/llm/health", nil)
	w := httptest.NewRecorder()

	controller.Execute(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.JSONEq(t, `{"status":"ok","llm":"up"}`, w.Body.String())
}

func TestLLMHealthCheckController_Error(t *testing.T) {
	mockClient := mockLLMClient{
		healthFn: func(ctx context.Context) error {
			return errors.New("llm down")
		},
	}
	handler := application.NewLLMHealthCheckHandler(mockClient)
	controller := controllers.NewLlmHealthCheckController(handler, zerolog.Nop())

	req := httptest.NewRequest(http.MethodGet, "/llm/health", nil)
	w := httptest.NewRecorder()

	controller.Execute(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Result().StatusCode)
	assert.JSONEq(t, `{"status":"degraded","llm":"down"}`, w.Body.String())
}

func TestLLMHealthCheckController_MethodNotAllowed(t *testing.T) {
	controller := controllers.NewLlmHealthCheckController(application.LLMHealthCheckHandler{}, zerolog.Nop())
	req := httptest.NewRequest(http.MethodPost, "/llm/health", nil)
	w := httptest.NewRecorder()

	controller.Execute(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Result().StatusCode)
}
