package controllers

import (
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
	w.WriteHeader(http.StatusOK)
}
