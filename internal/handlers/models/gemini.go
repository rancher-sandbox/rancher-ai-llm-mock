package modelHandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"llm-mock/internal/response"
	"llm-mock/internal/types"
)

type GeminiHandler struct {
	response *response.Handler
}

func NewGeminiHandler(response *response.Handler) *GeminiHandler {
	return &GeminiHandler{
		response: response,
	}
}

func (s *GeminiHandler) HandleRequest(c *gin.Context) {
	/**
	 * Expected path format: {model}:{api-name}
	 * example: gemini-flash-2.0:streamGenerateContent
	 */
	path := c.Param("path")
	parts := strings.Split(path, ":")

	switch parts[1] {
	case "streamGenerateContent":
		s.handleStreamGenerateContent(c)
	default:
		c.Status(404)
	}
}

func (s *GeminiHandler) handleStreamGenerateContent(c *gin.Context) {
	w := c.Writer
	alt := c.Query("alt")
	flusher, ok := w.(http.Flusher)
	if !ok {
		return
	}

	response := s.response.Pop()
	u := computeUsage(c, response)
	lastChunk := len(response.Text.Chunks) - 1

	if alt == "sse" {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// We accept only 1 MCP tool per response for simplicity
		if response.MCPTool.Name != "" {
			toolResp := s.buildToolResponse([]types.Tool{response.MCPTool}, u)
			b, _ := json.Marshal(toolResp)
			w.Write([]byte("data: "))
			w.Write(b)
			w.Write([]byte("\n\n"))
			flusher.Flush()
		} else if len(response.UITools) > 0 {
			toolResp := s.buildToolResponse(response.UITools, u)
			b, _ := json.Marshal(toolResp)
			w.Write([]byte("data: "))
			w.Write(b)
			w.Write([]byte("\n\n"))
			flusher.Flush()
		} else {
			for i, chunk := range response.Text.Chunks {
				resp := s.buildTextResponse(chunk, u, i == lastChunk)
				b, _ := json.Marshal(resp)
				w.Write([]byte("data: "))
				w.Write(b)
				w.Write([]byte("\n\n"))
				flusher.Flush()
				time.Sleep(100 * time.Millisecond)
			}
		}
		w.Write([]byte("data: {\"done\": true}\n\n"))
		flusher.Flush()
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Transfer-Encoding", "chunked")
		fmt.Fprint(w, "[")
		flusher.Flush()

		// We accept only 1 MCP tool per response for simplicity
		if response.MCPTool.Name != "" {
			toolResp := s.buildToolResponse([]types.Tool{response.MCPTool}, u)
			enc := json.NewEncoder(w)
			if err := enc.Encode(toolResp); err != nil {
				return
			}
		} else if len(response.UITools) > 0 {
			toolResp := s.buildToolResponse(response.UITools, u)
			enc := json.NewEncoder(w)
			if err := enc.Encode(toolResp); err != nil {
				return
			}
		} else {
			for i, chunk := range response.Text.Chunks {
				resp := s.buildTextResponse(chunk, u, i == lastChunk)
				enc := json.NewEncoder(w)
				if err := enc.Encode(resp); err != nil {
					return
				}
				if i < len(response.Text.Chunks)-1 {
					fmt.Fprint(w, ",")
				}
				flusher.Flush()
				time.Sleep(100 * time.Millisecond)
			}
		}

		// Close the array
		fmt.Fprint(w, "]")
		flusher.Flush()
	}
}

// Gemini usage counts are cumulative; sending them on the last chunk only gives the full totals.
func buildGeminiUsage(u usage) map[string]interface{} {
	return map[string]interface{}{
		"promptTokenCount":        u.InputTokens,
		"candidatesTokenCount":    u.OutputTokens,
		"totalTokenCount":         u.TotalTokens(),
		"cachedContentTokenCount": 0,
	}
}

func (s *GeminiHandler) buildTextResponse(chunk string, u usage, last bool) map[string]interface{} {
	resp := s.buildTextCandidates(chunk)
	if last {
		resp["usageMetadata"] = buildGeminiUsage(u)
	}
	return resp
}

func (s *GeminiHandler) buildTextCandidates(chunk string) map[string]interface{} {
	return map[string]interface{}{
		"candidates": []map[string]interface{}{
			{
				"content": map[string]interface{}{
					"parts": []map[string]interface{}{
						{"text": chunk},
					},
				},
				"finishReason": "STOP",
				"index":        0,
			},
		},
	}
}

func (s *GeminiHandler) buildToolResponse(tools []types.Tool, u usage) map[string]interface{} {
	parts := make([]map[string]interface{}, len(tools))

	for i, tool := range tools {
		parts[i] = map[string]interface{}{
			"function_call": map[string]interface{}{
				"name": tool.Name,
				"args": tool.Args,
			},
		}
	}

	return map[string]interface{}{
		"candidates": []map[string]interface{}{
			{
				"content": map[string]interface{}{
					"parts": parts,
				},
				"finishReason": "STOP",
				"index":        0,
			},
		},
		"usageMetadata": buildGeminiUsage(u),
	}
}
