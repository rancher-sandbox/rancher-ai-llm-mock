package internal

import (
	control "rancher-ai-llm-mock/internal/handlers"
	handlers "rancher-ai-llm-mock/internal/handlers/models"
	"rancher-ai-llm-mock/internal/queue"
)

type modelHandlers struct {
	Gemini *handlers.GeminiHandler
	OpenAI *handlers.OpenAIHandler
}

type llmService struct {
	queue   *queue.Queue
	Control *control.ControlHandler
	Models  modelHandlers
}

func NewLLMService() *llmService {
	queue := queue.NewQueue()

	return &llmService{
		queue: queue,
		Models: modelHandlers{
			Gemini: handlers.NewGeminiHandler(queue),
			OpenAI: handlers.NewOpenAIHandler(queue),
		},
		Control: control.NewControlHandler(queue),
	}
}
