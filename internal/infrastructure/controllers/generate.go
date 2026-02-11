package controllers

import (
	"encoding/json"
	"llm-agent-go/internal/application"
	"net/http"
)

type GenerateController struct {
	handler application.GenerateHandler
}

type GenerateRequest struct {
	Prompt string `json:"prompt"`
}

func NewGenerateController(
	handler application.GenerateHandler,
) GenerateController {
	return GenerateController{
		handler: handler,
	}
}

func (c GenerateController) Execute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	if req.Prompt == "" {
		http.Error(w, "prompt required", http.StatusBadRequest)
		return
	}

	resp, err := c.handler.Handle(r.Context(), req.Prompt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})

		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"response": resp,
	})
}
