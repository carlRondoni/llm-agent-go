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
	w.WriteHeader(http.StatusOK)
}
