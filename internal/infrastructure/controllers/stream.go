package controllers

import (
	"encoding/json"
	"llm-agent-go/internal/application"
	"net/http"

	"github.com/rs/zerolog"
)

type StreamController struct {
	handler application.StreamHandler
	logger  zerolog.Logger
}

type StreamRequest struct {
	Prompt string `json:"prompt"`
}

func NewStreamController(
	handler application.StreamHandler,
	logger zerolog.Logger,
) StreamController {
	return StreamController{
		handler: handler,
		logger:  logger,
	}
}

func (c StreamController) Execute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		c.logger.Error().Msg("stream error: method not allowed")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	flusher, ok := w.(http.Flusher)
	if !ok {
		c.logger.Error().Msg("stream error: flusher unsupported")
		http.Error(w, "stream unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Transfer-Encoding", "chunked")

	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.logger.Error().Err(err).Msg("stream error: invalid json body")
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	if req.Prompt == "" {
		c.logger.Error().Msg("stream error: prompt required")
		http.Error(w, "prompt required", http.StatusBadRequest)
		return
	}

	stream, err := c.handler.Handle(ctx, req.Prompt)
	if err != nil {
		c.logger.Error().Err(err).Msg("stream error: error on handler: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for token := range stream {
		_, _ = w.Write([]byte(token))
		flusher.Flush()
	}
}
