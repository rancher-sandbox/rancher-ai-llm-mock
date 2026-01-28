package internal

import (
	control "llm-mock/internal/handlers"
	handlers "llm-mock/internal/handlers/models"
	response "llm-mock/internal/response"
)

type modelHandlers struct {
	Ollama  *handlers.OllamaHandler
	Gemini  *handlers.GeminiHandler
	OpenAI  *handlers.OpenAIHandler
	Bedrock *handlers.BedrockHandler
}

type llmService struct {
	responseHandler *response.Handler
	Control         *control.ControlHandler
	Models          modelHandlers
}

func NewLLMService() *llmService {
	responseHandler := response.NewHandler()

	return &llmService{
		responseHandler: responseHandler,
		Models: modelHandlers{
			Gemini:  handlers.NewGeminiHandler(responseHandler),
			OpenAI:  handlers.NewOpenAIHandler(responseHandler),
			Ollama:  handlers.NewOllamaHandler(responseHandler),
			Bedrock: handlers.NewBedrockHandler(responseHandler),
		},
		Control: control.NewControlHandler(responseHandler),
	}
}
