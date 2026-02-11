package service_container

type ServiceContainer struct {
	Controllers
}

func NewServiceContainer() ServiceContainer {
	llmClients := NewLLMClients()
	handlers := NewHandlers(llmClients)
	controllers := NewControllers(handlers)

	return ServiceContainer{
		Controllers: controllers,
	}
}
