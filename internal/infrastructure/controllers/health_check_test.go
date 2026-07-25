package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"llm-agent-go/internal/infrastructure/controllers"
)

func TestHealthCheckController_Execute_Success(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	controller := controllers.NewHealthCheckController(zerolog.Nop())
	controller.Execute(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}

func TestHealthCheckController_Execute_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	w := httptest.NewRecorder()

	controller := controllers.NewHealthCheckController(zerolog.Nop())
	controller.Execute(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Result().StatusCode)
}
