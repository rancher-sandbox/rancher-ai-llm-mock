package main

import (
	"fmt"
	"log"
	"rancher-ai-llm-mock/internal"

	"github.com/gin-gonic/gin"
)

func main() {
	port := ":8083"

	log.Println(fmt.Sprintf("Starting LLM Mock Service on %s...", port))

	r := gin.Default()
	svc := internal.NewLLMService()

	// Control endpoints
	r.POST("/v1/control/push", svc.Control.HandlePushRequest)

	r.POST("/v1/control/clear", svc.Control.HandleClearRequest)

	// Ollama model endpoint
	r.POST("/api/chat", svc.Models.Ollama.HandleRequest)

	// Gemini model endpoint
	r.POST("/v1beta/models/:path", svc.Models.Gemini.HandleRequest)

	// OpenAI model endpoint
	r.POST("/chat/completions", svc.Models.OpenAI.HandleRequest)

	r.Run(port)
}
