package controllers

import (
	"encoding/json"
	"llm-agent-go/internal/application"
	"net/http"
)

type LlmHealthCheckController struct {
	handler application.LLMHealthCheckHandler
}

func NewLlmHealthCheckController(
	handler application.LLMHealthCheckHandler,
) LlmHealthCheckController {
	return LlmHealthCheckController{
		handler: handler,
	}
}

func (c LlmHealthCheckController) Execute(w http.ResponseWriter, r *http.Request) {
	err := c.handler.Handle(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]any{
			"status": "degraded",
			"llm":    "down",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"status": "ok",
		"llm":    "up",
	})
}
