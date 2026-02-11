package controllers

import (
	"llm-agent-go/internal/application"
	"net/http"
)

type QueryController struct {
	handler application.QueryHandler
}

func NewQueryController(
	handler application.QueryHandler,
) QueryController {
	return QueryController{
		handler: handler,
	}
}

func (c QueryController) Execute(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
