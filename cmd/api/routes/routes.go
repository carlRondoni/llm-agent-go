package routes

import (
	"llm-agent-go/cmd/service_container"
	"net/http"
)

func InitRoutes(controllers service_container.Controllers) {
	// basic routes
	http.Handle("/health", http.HandlerFunc(controllers.HealthCheckController.Execute))

	// llm routes
	http.Handle("/llm/generate", http.HandlerFunc(controllers.GenerateController.Execute))
	http.Handle("/llm/stream", http.HandlerFunc(controllers.StreamController.Execute))
	http.Handle("/llm/health", http.HandlerFunc(controllers.LlmHealthCheckController.Execute))
}
