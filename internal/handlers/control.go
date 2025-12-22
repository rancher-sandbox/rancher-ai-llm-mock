package control

import (
	"rancher-ai-llm-mock/internal/queue"
	"rancher-ai-llm-mock/internal/types"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	queue *queue.Queue
}

func NewHandler(queue *queue.Queue) *Handler {
	return &Handler{
		queue: queue,
	}
}

func (s *Handler) HandlePushRequest(c *gin.Context) {
	var req types.MockResponse
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(400)
		return
	}
	s.queue.Push(req.Chunks)
	c.Status(204)
}

func (s *Handler) HandleClearRequest(c *gin.Context) {
	s.queue.Clear()
	c.Status(204)
}
