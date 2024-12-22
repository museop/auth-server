package store

import (
	"fmt"
	"sync"
)

// InMemoryUserStore는 메모리 기반 UserStore 구현체
type InMemoryUserStore struct {
	store map[string]string
	mutex sync.Mutex
}

// NewInMemoryUserStore는 새로운 InMemoryUserStore 인스턴스를 생성합니다.
func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		store: make(map[string]string),
	}
}

// SaveUser는 사용자 이름과 해시된 비밀번호를 저장합니다.
func (s *InMemoryUserStore) SaveUser(username, hashedPassword string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 중복 사용자 체크
	if _, exists := s.store[username]; exists {
		return fmt.Errorf("user already exists")
	}

	s.store[username] = hashedPassword
	return nil
}

// GetUser는 사용자 이름으로 해시된 비빌번호를 가져옵니다.
func (s *InMemoryUserStore) GetUser(username string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	hashedPassword, exists := s.store[username]
	if !exists {
		return "", fmt.Errorf("user not found")
	}

	return hashedPassword, nil
}
