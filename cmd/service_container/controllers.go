package service_container

import (
	"llm-agent-go/internal/infrastructure/controllers"

	"github.com/rs/zerolog"
)

type Controllers struct {
	HealthCheckController    controllers.HealthCheckController
	LlmHealthCheckController controllers.LlmHealthCheckController
	GenerateController       controllers.GenerateController
	StreamController         controllers.StreamController
}

func NewControllers(handlers Handlers, logger zerolog.Logger) Controllers {
	return Controllers{
		HealthCheckController:    controllers.NewHealthCheckController(logger),
		LlmHealthCheckController: controllers.NewLlmHealthCheckController(handlers.LLMHealthHandler, logger),
		GenerateController:       controllers.NewGenerateController(handlers.GenerateHandler, logger),
		StreamController:         controllers.NewStreamController(handlers.StreamHandler, logger),
	}
}
