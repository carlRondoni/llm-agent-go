package controllers

import (
	"net/http"
)

type HealthCheckController struct{}

func NewHealthCheckController() HealthCheckController {
	return HealthCheckController{}
}

func (c HealthCheckController) Execute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	w.WriteHeader(http.StatusOK)
}
