package modelHandlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"llm-mock/internal/response"
	"llm-mock/internal/types"
)

type OpenAIHandler struct {
	response *response.Handler
}

func NewOpenAIHandler(response *response.Handler) *OpenAIHandler {
	return &OpenAIHandler{
		response: response,
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

	response := s.response.Pop()

	if response.Tool.Name != "" {
		b := s.buildToolResponse([]types.Tool{response.Tool})
		w.Write([]byte("data: "))
		w.Write(b)
		w.Write([]byte("\n\n"))
		flusher.Flush()
	} else if len(response.UITools) > 0 {
		b := s.buildToolResponse(response.UITools)
		w.Write([]byte("data: "))
		w.Write(b)
		w.Write([]byte("\n\n"))
		flusher.Flush()
	} else {
		for i, chunk := range response.Text.Chunks {
			b := s.buildTextResponse(chunk, i)
			w.Write([]byte("data: "))
			w.Write(b)
			w.Write([]byte("\n\n"))
			flusher.Flush()
			time.Sleep(100 * time.Millisecond)
		}
	}

	b := s.buildFinalChunk()
	w.Write([]byte("data: "))
	w.Write(b)
	w.Write([]byte("\n\n"))
	flusher.Flush()

	// End the stream
	w.Write([]byte("data: [DONE]\n\n"))
	flusher.Flush()
}

func (s *OpenAIHandler) buildTextResponse(chunk string, i int) []byte {
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
				"finish_reason": "stop",
			},
		},
	}
	b, _ := json.Marshal(data)

	return b
}

func (s *OpenAIHandler) buildToolResponse(tools []types.Tool) []byte {
	toolCalls := make([]map[string]interface{}, len(tools))

	for i, tool := range tools {
		argsStr := "{}"
		if tool.Args != nil {
			if b, err := json.Marshal(tool.Args); err == nil {
				argsStr = string(b)
			}
		}

		toolCalls[i] = map[string]interface{}{
			"id":   "call_abc123",
			"type": "function",
			"function": map[string]interface{}{
				"name":      tool.Name,
				"arguments": argsStr,
			},
		}
	}

	data := map[string]interface{}{
		"id":      "chatcmpl-mock-12345",
		"object":  "chat.completion.chunk",
		"created": time.Now().Unix(),
		"model":   "openai-mock-v1",
		"choices": []map[string]interface{}{
			{
				"index": 0,
				"delta": map[string]interface{}{
					"role":       "assistant",
					"tool_calls": toolCalls,
				},
				"finish_reason": "tool_calls",
			},
		},
	}
	b, _ := json.Marshal(data)

	return b
}

func (s *OpenAIHandler) buildFinalChunk() []byte {
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

	return b
}
