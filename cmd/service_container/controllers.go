package service_container

import "llm-agent-go/internal/infrastructure/controllers"

type Controllers struct {
	HealthCheckController    controllers.HealthCheckController
	LlmHealthCheckController controllers.LlmHealthCheckController
	QueryController          controllers.QueryController
	StreamController         controllers.StreamController
}

func NewControllers(handlers Handlers) Controllers {
	return Controllers{
		HealthCheckController:    controllers.NewHealthCheckController(),
		LlmHealthCheckController: controllers.NewLlmHealthCheckController(handlers.LLMHealthHandler),
		QueryController:          controllers.NewQueryController(handlers.QueryHandler),
		StreamController:         controllers.NewStreamController(handlers.StreamHandler),
	}
}
