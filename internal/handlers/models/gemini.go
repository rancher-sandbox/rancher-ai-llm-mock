package gemini

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"rancher-ai-llm-mock/internal/queue"
	types "rancher-ai-llm-mock/internal/types"
)

type Handler struct {
	queue *queue.Queue
}

func NewHandler(queue *queue.Queue) *Handler {
	return &Handler{
		queue: queue,
	}
}

func (s *Handler) HandleRequest(c *gin.Context) {
	path := c.Param("path")

	// "path" is {model}:{some-gemini-api-name}
	parts := strings.Split(path, ":")

	switch parts[1] {
	case "streamGenerateContent":
		s.HandleStreamGenerateContent(c)
	default:
		c.Status(404)
	}
}

func (s *Handler) HandleStreamGenerateContent(c *gin.Context) {
	w := c.Writer
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Transfer-Encoding", "chunked")

	flusher, ok := w.(http.Flusher)
	if !ok {
		return
	}

	fmt.Fprint(w, "[")
	flusher.Flush()

	response := s.queue.Pop()

	var chunks []string
	if len(response.Chunks) == 0 {
		chunks = []string{"Mock service queue is empty.", " This is", " a default mock response", " from the Gemini model."}
	} else {
		chunks = response.Chunks
	}

	for i, text := range chunks {
		resp := types.GenerateContentResponse{
			Candidates: []types.Candidate{
				{
					Content: types.Content{
						Parts: []types.Part{
							{Text: text},
						},
					},
					FinishReason: "length",
					Index:        0,
				},
			},
			ModelVersion: "gemini-mock-v1",
			ResponseId:   "resp-mock-12345",
		}

		// Use the encoder directly on the writer
		enc := json.NewEncoder(w)
		if err := enc.Encode(resp); err != nil {
			return
		}

		// Handle commas
		if i < len(chunks)-1 {
			fmt.Fprint(w, ",")
		}

		flusher.Flush()
		time.Sleep(200 * time.Millisecond)
	}

	// Close the array
	fmt.Fprint(w, "]")
	flusher.Flush()
}
