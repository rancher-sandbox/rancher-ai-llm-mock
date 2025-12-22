package queue

import (
	types "rancher-ai-llm-mock/internal/types"
	"sync"
)

type Queue struct {
	mu       sync.RWMutex
	messages []types.MockResponse
}

func NewQueue() *Queue {
	return &Queue{
		messages: []types.MockResponse{},
	}
}

func (q *Queue) Push(chunks []string) {
	q.mu.Lock()
	q.messages = append(q.messages, types.MockResponse{Chunks: chunks})
	q.mu.Unlock()
}

func (q *Queue) Pop() types.MockResponse {
	q.mu.RLock()
	if len(q.messages) == 0 {
		q.mu.RUnlock()
		return types.MockResponse{Chunks: []string{}}
	}

	resp := q.messages[0]
	q.messages = q.messages[1:]
	q.mu.RUnlock()
	return resp
}

func (q *Queue) Clear() {
	q.mu.Lock()
	q.messages = []types.MockResponse{}
	q.mu.Unlock()
}
