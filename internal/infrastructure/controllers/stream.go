package controllers

import (
	"encoding/json"
	"llm-agent-go/internal/application"
	"net/http"
)

type StreamController struct {
	handler application.StreamHandler
}

type StreamRequest struct {
	Prompt string `json:"prompt"`
}

func NewStreamController(
	handler application.StreamHandler,
) StreamController {
	return StreamController{
		handler: handler,
	}
}

func (c StreamController) Execute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "stream unsupported", 500)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Transfer-Encoding", "chunked")

	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	if req.Prompt == "" {
		http.Error(w, "prompt required", http.StatusBadRequest)
		return
	}

	stream, err := c.handler.Handle(ctx, req.Prompt)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	for token := range stream {
		_, _ = w.Write([]byte(token))
		flusher.Flush()
	}
}
