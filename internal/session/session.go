package session

import (
	"sync"
)

type Session struct {
	mu    sync.RWMutex
	token string
}

func NewSession() *Session {
	return &Session{}
}

func (s *Session) SetToken(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.token = token
}

func (s *Session) GetToken() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.token
}
