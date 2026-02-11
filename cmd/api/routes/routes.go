package routes

import (
	"llm-agent-go/cmd/service_container"
	"net/http"
)

func InitRoutes(controllers service_container.Controllers) {
	// health checks
	http.Handle("/health", http.HandlerFunc(controllers.HealthCheckController.Execute))
	http.Handle("/health/llm", http.HandlerFunc(controllers.LlmHealthCheckController.Execute))

	// llm endpoints
	http.Handle("/llm/generate", http.HandlerFunc(controllers.GenerateController.Execute))
	http.Handle("/llm/stream", http.HandlerFunc(controllers.StreamController.Execute))
}
