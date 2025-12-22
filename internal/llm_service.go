package internal

import (
	control "rancher-ai-llm-mock/internal/handlers"
	gemini "rancher-ai-llm-mock/internal/handlers/models"
	"rancher-ai-llm-mock/internal/queue"
)

type modelHandlers struct {
	Gemini *gemini.Handler
}

type llmService struct {
	queue   *queue.Queue
	Control *control.Handler
	Models  modelHandlers
}

func NewLLMService() *llmService {
	queue := queue.NewQueue()

	return &llmService{
		queue: queue,
		Models: modelHandlers{
			Gemini: gemini.NewHandler(queue),
		},
		Control: control.NewHandler(queue),
	}
}
