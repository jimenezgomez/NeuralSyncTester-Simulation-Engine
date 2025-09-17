package session_manager

import (
	"sync"
	"time"

	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine"
)

// SSESession represents a client session.
type SSESession struct {
	UID            string
	IP             string
	SyncSessionUID string
	Expiry         time.Time
}

// SessionManager manages active sessions.
type SessionManager struct {
	sessionsSSE  map[string]SSESession
	sessionsMTPM map[string]*engine.TrackedMTPMState
	mu           sync.RWMutex
	ttl          time.Duration
}

// NewSessionManager creates a manager with given session TTL.
func NewSessionManager(ttl time.Duration) *SessionManager {
	sm := &SessionManager{
		sessionsSSE:  make(map[string]SSESession),
		sessionsMTPM: make(map[string]*engine.TrackedMTPMState),
		ttl:          ttl,
	}
	// start background cleanup loop
	go sm.updateLoop()
	return sm
}

// AddSSE adds or updates a session.
func (sm *SessionManager) AddSSE(token string, s SSESession) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	s.Expiry = time.Now().Add(sm.ttl)
	sm.sessionsSSE[token] = s
}

// GetSSE retrieves a session by token, if valid.
func (sm *SessionManager) GetSSE(token string) (SSESession, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	s, ok := sm.sessionsSSE[token]
	if !ok {
		return SSESession{}, false
	}
	if time.Now().After(s.Expiry) {
		return SSESession{}, false
	}
	return s, true
}

// DeleteSSE removes a session manually.
func (sm *SessionManager) DeleteSSE(token string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessionsSSE, token)
}

// AddSSE adds or updates a session.
func (sm *SessionManager) AddMTPM(token string, s *engine.TrackedMTPMState) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.sessionsMTPM[token] = s
}

// GetSSE retrieves a session by token, if valid.
func (sm *SessionManager) GetMTPM(token string) (*engine.TrackedMTPMState, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	s, ok := sm.sessionsMTPM[token]
	if !ok {
		return nil, false
	}
	return s, true
}

// DeleteSSE removes a session manually.
func (sm *SessionManager) DeleteMTPM(token string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessionsMTPM, token)
}

// updateLoop runs every 10 seconds to remove expired sessions.
func (sm *SessionManager) updateLoop() {
	tickerSubs := time.NewTicker(1 * time.Second)
	defer tickerSubs.Stop()
	tickerExpired := time.NewTicker(1 * time.Minute)
	defer tickerExpired.Stop()
	for {
		select {
		case <-tickerSubs.C:
			sm.updateAllSnapshots()
		case <-tickerExpired.C:
			sm.cleanupExpired()
		}
	}
}

// cleanupExpired removes expired sessions.
func (sm *SessionManager) cleanupExpired() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	now := time.Now()
	for token, s := range sm.sessionsSSE {
		if now.After(s.Expiry) || s.Expiry.IsZero() {
			delete(sm.sessionsSSE, token)
		}
	}
}

func (sm *SessionManager) updateAllSnapshots() {
	for _, v := range sm.sessionsMTPM {
		v.UpdateAllSubscribers()
	}
}
