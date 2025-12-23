package modelHandlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"rancher-ai-llm-mock/internal/queue"
)

type OllamaHandler struct {
	queue *queue.Queue
}

func NewOllamaHandler(queue *queue.Queue) *OllamaHandler {
	return &OllamaHandler{
		queue: queue,
	}
}

func (s *OllamaHandler) HandleRequest(c *gin.Context) {
	w := c.Writer
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Transfer-Encoding", "chunked")

	flusher, ok := w.(http.Flusher)
	if !ok {
		return
	}

	response := s.queue.Pop()

	for i, text := range response.Chunks {
		resp := map[string]interface{}{
			"model":      "ollama-mock-v1",
			"created_at": time.Now().Format(time.RFC3339),
			"message": map[string]interface{}{
				"role":    "assistant",
				"content": text,
			},
			"done": i == len(response.Chunks)-1,
		}
		enc := json.NewEncoder(w)
		if err := enc.Encode(resp); err != nil {
			return
		}
		flusher.Flush()
		time.Sleep(200 * time.Millisecond)
	}
}
