package service

import (
	"sync"
	"time"

	"github.com/spaghetti-people/khuthon-ai/internal/model"
)

// ConversationStore 대화 기록을 저장하는 인터페이스
type ConversationStore interface {
	// SaveConversation 대화 기록을 저장합니다
	SaveConversation(plantID string, messages []model.Message) error
	// GetConversation 저장된 대화 기록을 조회합니다
	GetConversation(plantID string) ([]model.Message, error)
}

// MemoryConversationStore 메모리에 대화 기록을 저장하는 구현체
type MemoryConversationStore struct {
	mu            sync.RWMutex
	conversations map[string][]model.Message
}

// NewMemoryConversationStore 새로운 메모리 대화 저장소를 생성합니다
func NewMemoryConversationStore() *MemoryConversationStore {
	return &MemoryConversationStore{
		conversations: make(map[string][]model.Message),
	}
}

// SaveConversation 대화 기록을 저장합니다
func (s *MemoryConversationStore) SaveConversation(plantID string, messages []model.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 현재 시간을 포함한 메시지 복사본 생성
	now := time.Now().Format(time.RFC3339)
	copiedMessages := make([]model.Message, len(messages))
	for i, msg := range messages {
		copiedMessages[i] = model.Message{
			Role:    msg.Role,
			Content: msg.Content,
			Time:    now,
		}
	}

	s.conversations[plantID] = copiedMessages
	return nil
}

// GetConversation 저장된 대화 기록을 조회합니다
func (s *MemoryConversationStore) GetConversation(plantID string) ([]model.Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if messages, exists := s.conversations[plantID]; exists {
		return messages, nil
	}
	return []model.Message{}, nil
}
