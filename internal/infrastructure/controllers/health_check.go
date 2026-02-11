package controllers

import (
	"net/http"
)

type HealthCheckController struct{}

func NewHealthCheckController() HealthCheckController {
	return HealthCheckController{}
}

func (c HealthCheckController) Execute(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
