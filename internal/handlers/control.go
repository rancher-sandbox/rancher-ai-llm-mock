package controlHandler

import (
	"llm-mock/internal/response"
	"llm-mock/internal/types"

	"github.com/gin-gonic/gin"
)

type ControlHandler struct {
	response *response.Handler
}

func NewControlHandler(response *response.Handler) *ControlHandler {
	return &ControlHandler{
		response: response,
	}
}

func (s *ControlHandler) HandlePushRequest(c *gin.Context) {
	var req types.MockResponse

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(400)
		return
	}

	if req.MCPTool.Name != "" && req.MCPTool.Args == nil {
		c.JSON(400, gin.H{"error": "Invalid payload: MCPTool.Args must be provided when MCPTool is set"})
		return
	}

	if req.MCPTool.Name == "" && req.MCPTool.Args != nil {
		c.JSON(400, gin.H{"error": "Invalid payload: MCPTool.Name must be provided when MCPTool is set"})
		return
	}

	if (req.Text.Chunks == nil || len(req.Text.Chunks) == 0) && (req.MCPTool.Name == "" || req.MCPTool.Args == nil) {
		c.JSON(400, gin.H{"error": "Invalid payload: one of Text or MCPTool fields must be provided"})
		return
	}

	s.response.Push(req)

	c.Status(204)
}

func (s *ControlHandler) HandleClearRequest(c *gin.Context) {
	s.response.Clear()
	c.Status(204)
}
