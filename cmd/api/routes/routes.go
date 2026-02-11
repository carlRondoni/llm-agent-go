package routes

import (
	"llm-agent-go/cmd/service_container"
	"net/http"
)

func InitRoutes(controllers service_container.Controllers) {
	// health checks
	http.Handle("/health", http.HandlerFunc(controllers.HealthCheckController.Execute))
	http.Handle("/v1/llm/health", http.HandlerFunc(controllers.LlmHealthCheckController.Execute))

	// llm endpoints
	http.Handle("/v1/llm/query", http.HandlerFunc(controllers.QueryController.Execute))
	http.Handle("/v1/llm/stream", http.HandlerFunc(controllers.StreamController.Execute))
}
