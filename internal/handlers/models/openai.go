package modelHandlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"llm-mock/internal/queue"
)

type OpenAIHandler struct {
	queue *queue.Queue
}

func NewOpenAIHandler(queue *queue.Queue) *OpenAIHandler {
	return &OpenAIHandler{
		queue: queue,
	}
}

func (s *OpenAIHandler) HandleRequest(c *gin.Context) {
	w := c.Writer
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	response := s.queue.Pop()

	for i, chunk := range response.Chunks {
		data := map[string]interface{}{
			"id":      "chatcmpl-mock-12345",
			"object":  "chat.completion.chunk",
			"created": time.Now().Unix(),
			"model":   "openai-mock-v1",
			"choices": []map[string]interface{}{
				{
					"index": 0,
					"delta": map[string]interface{}{
						"role": func() interface{} {
							if i == 0 {
								return "assistant"
							} else {
								return nil
							}
						}(),
						"content": chunk,
					},
					"finish_reason": nil,
				},
			},
		}
		b, _ := json.Marshal(data)
		w.Write([]byte("data: "))
		w.Write(b)
		w.Write([]byte("\n\n"))
		flusher.Flush()
		time.Sleep(200 * time.Millisecond)
	}

	// Send the final chunk with finish_reason
	data := map[string]interface{}{
		"id":      "chatcmpl-mock-12345",
		"object":  "chat.completion.chunk",
		"created": time.Now().Unix(),
		"model":   "openai-mock-v1",
		"choices": []map[string]interface{}{
			{
				"index":         0,
				"delta":         map[string]interface{}{},
				"finish_reason": "stop",
			},
		},
	}
	b, _ := json.Marshal(data)
	w.Write([]byte("data: "))
	w.Write(b)
	w.Write([]byte("\n\n"))
	flusher.Flush()

	// End the stream
	w.Write([]byte("data: [DONE]\n\n"))
	flusher.Flush()
}
