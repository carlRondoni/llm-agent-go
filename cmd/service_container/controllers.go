package service_container

import "llm-agent-go/internal/infrastructure/controllers"

type Controllers struct {
	HealthCheckController    controllers.HealthCheckController
	LlmHealthCheckController controllers.LlmHealthCheckController
	GenerateController       controllers.GenerateController
	StreamController         controllers.StreamController
}

func NewControllers(handlers Handlers) Controllers {
	return Controllers{
		HealthCheckController:    controllers.NewHealthCheckController(),
		LlmHealthCheckController: controllers.NewLlmHealthCheckController(handlers.LLMHealthHandler),
		GenerateController:       controllers.NewGenerateController(handlers.GenerateHandler),
		StreamController:         controllers.NewStreamController(handlers.StreamHandler),
	}
}
