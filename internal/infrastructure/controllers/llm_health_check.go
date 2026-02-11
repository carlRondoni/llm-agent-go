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
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
