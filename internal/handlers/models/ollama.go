package modelHandlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"llm-mock/internal/response"
	"llm-mock/internal/types"
)

type OllamaHandler struct {
	response *response.Handler
}

func NewOllamaHandler(response *response.Handler) *OllamaHandler {
	return &OllamaHandler{
		response: response,
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

	response := s.response.Pop()
	u := computeUsage(c, response)

	if response.MCPTool.Name != "" {
		resp := s.buildToolResponse([]types.Tool{response.MCPTool}, u)
		enc := json.NewEncoder(w)
		if err := enc.Encode(resp); err != nil {
			return
		}
		flusher.Flush()
	} else if len(response.UITools) > 0 {
		resp := s.buildToolResponse(response.UITools, u)
		enc := json.NewEncoder(w)
		if err := enc.Encode(resp); err != nil {
			return
		}
		flusher.Flush()
	} else {
		for i, text := range response.Text.Chunks {
			resp := s.buildTextResponse(text, i == len(response.Text.Chunks)-1, u)
			enc := json.NewEncoder(w)
			if err := enc.Encode(resp); err != nil {
				return
			}
			flusher.Flush()
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (s *OllamaHandler) buildTextResponse(chunk string, done bool, u usage) map[string]interface{} {
	resp := map[string]interface{}{
		"message": map[string]interface{}{
			"role":    "assistant",
			"content": chunk,
		},
		"done": done,
	}
	if done {
		resp["done_reason"] = "stop"
		addOllamaUsage(resp, u)
	}
	return resp
}

// Ollama reports token counts only on the final (done) chunk.
func addOllamaUsage(resp map[string]interface{}, u usage) {
	resp["prompt_eval_count"] = u.InputTokens
	resp["eval_count"] = u.OutputTokens
}

func (s *OllamaHandler) buildToolResponse(tools []types.Tool, u usage) map[string]interface{} {
	toolCalls := make([]map[string]interface{}, len(tools))

	for i, tool := range tools {
		toolCalls[i] = map[string]interface{}{
			"function": map[string]interface{}{
				"name":      tool.Name,
				"arguments": tool.Args,
			},
		}
	}

	resp := map[string]interface{}{
		"message": map[string]interface{}{
			"role":       "assistant",
			"content":    "",
			"tool_calls": toolCalls,
		},
		"done_reason": "stop",
		"done":        true,
	}
	addOllamaUsage(resp, u)
	return resp
}
