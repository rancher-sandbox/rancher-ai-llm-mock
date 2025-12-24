package internal

import (
	control "llm-mock/internal/handlers"
	handlers "llm-mock/internal/handlers/models"
	"llm-mock/internal/queue"
)

type modelHandlers struct {
	Ollama  *handlers.OllamaHandler
	Gemini  *handlers.GeminiHandler
	OpenAI  *handlers.OpenAIHandler
	Bedrock *handlers.BedrockHandler
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
			Gemini:  handlers.NewGeminiHandler(queue),
			OpenAI:  handlers.NewOpenAIHandler(queue),
			Ollama:  handlers.NewOllamaHandler(queue),
			Bedrock: handlers.NewBedrockHandler(queue),
		},
		Control: control.NewControlHandler(queue),
	}
}
