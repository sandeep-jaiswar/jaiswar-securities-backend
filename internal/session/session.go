package session

import (
	"sync"
)

type SessionManager struct {
	tokens map[string]string
	mu     sync.Mutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		tokens: make(map[string]string),
	}
}

func (sm *SessionManager) StoreToken(userID, token string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.tokens[userID] = token
}

func (sm *SessionManager) GetToken(userID string) (string, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	token, exists := sm.tokens[userID]
	return token, exists
}

func (sm *SessionManager) DeleteToken(userID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.tokens, userID)
}
