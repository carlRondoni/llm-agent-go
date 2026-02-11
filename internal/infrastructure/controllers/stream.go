package controllers

import (
	"llm-agent-go/internal/application"
	"net/http"
)

type StreamController struct {
	handler application.StreamHandler
}

func NewStreamController(
	handler application.StreamHandler,
) StreamController {
	return StreamController{
		handler: handler,
	}
}

func (c StreamController) Execute(w http.ResponseWriter, r *http.Request) {
ctx := r.Context()

    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "stream unsupported", 500)
        return
    }

    w.Header().Set("Content-Type", "text/plain")
    w.Header().Set("Transfer-Encoding", "chunked")

    prompt := r.URL.Query().Get("q")

    stream, err := c.handler.Handle(ctx, prompt)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    for token := range stream {
        _, _ = w.Write([]byte(token))
        flusher.Flush()
    }
}
