package modelHandlers

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/gin-gonic/gin"

	"llm-mock/internal/types"
)

// usage is a deterministic token count estimate, using the ~4 characters per token rule of thumb.
type usage struct {
	InputTokens  int
	OutputTokens int
}

func (u usage) TotalTokens() int {
	return u.InputTokens + u.OutputTokens
}

func estimateTokens(s string) int {
	if s == "" {
		return 0
	}
	return (len(s) + 3) / 4
}

func computeUsage(c *gin.Context, resp types.MockResponse) usage {
	body, _ := io.ReadAll(c.Request.Body)

	var out strings.Builder
	out.WriteString(strings.Join(resp.Text.Chunks, ""))
	tools := resp.UITools
	if resp.MCPTool.Name != "" {
		tools = append([]types.Tool{resp.MCPTool}, tools...)
	}
	for _, tool := range tools {
		out.WriteString(tool.Name)
		if b, err := json.Marshal(tool.Args); err == nil {
			out.Write(b)
		}
	}

	return usage{
		InputTokens:  max(estimateTokens(string(body)), 1),
		OutputTokens: max(estimateTokens(out.String()), 1),
	}
}
