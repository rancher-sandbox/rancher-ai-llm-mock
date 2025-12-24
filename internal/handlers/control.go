package controlHandler

import (
	"llm-mock/internal/queue"
	"llm-mock/internal/types"

	"github.com/gin-gonic/gin"
)

type ControlHandler struct {
	queue *queue.Queue
}

func NewControlHandler(queue *queue.Queue) *ControlHandler {
	return &ControlHandler{
		queue: queue,
	}
}

func (s *ControlHandler) HandlePushRequest(c *gin.Context) {
	var req types.MockResponse
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(400)
		return
	}
	s.queue.Push(req.Chunks)
	c.Status(204)
}

func (s *ControlHandler) HandleClearRequest(c *gin.Context) {
	s.queue.Clear()
	c.Status(204)
}
