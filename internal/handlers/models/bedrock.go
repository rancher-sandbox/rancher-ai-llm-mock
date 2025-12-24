package modelHandlers

import (
	"github.com/gin-gonic/gin"

	"llm-mock/internal/queue"
)

type BedrockHandler struct {
	queue *queue.Queue
}

func NewBedrockHandler(queue *queue.Queue) *BedrockHandler {
	return &BedrockHandler{
		queue: queue,
	}
}

func (s *BedrockHandler) HandleRequest(c *gin.Context) {
	// TODO: Implement Bedrock-specific request handling logic here
}
