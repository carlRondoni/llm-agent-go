package controllers

import (
	"encoding/json"
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
	prompt := r.URL.Query().Get("q")

	resp, err := c.handler.Handle(r.Context(), prompt)
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
