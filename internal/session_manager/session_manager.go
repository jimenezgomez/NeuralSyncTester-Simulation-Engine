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
	sessionsMTPM map[string]*engine.TrackedMTPMSession
	mu           sync.RWMutex
	ttl          time.Duration
}

// NewSessionManager creates a manager with given session TTL.
func NewSessionManager(ttl time.Duration) *SessionManager {
	sm := &SessionManager{
		sessionsSSE:  make(map[string]SSESession),
		sessionsMTPM: make(map[string]*engine.TrackedMTPMSession),
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
func (sm *SessionManager) AddMTPM(token string, s *engine.TrackedMTPMSession) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.sessionsMTPM[token] = s
}

// GetSSE retrieves a session by token, if valid.
func (sm *SessionManager) GetMTPM(token string) (*engine.TrackedMTPMSession, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	s, ok := sm.sessionsMTPM[token]
	if !ok {
		return nil, false
	}
	return s, true
}

func (sm *SessionManager) GetAllMTPMSessions() map[string]*engine.TrackedMTPMSession {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// Make a shallow copy to avoid concurrent map access issues
	copyMap := make(map[string]*engine.TrackedMTPMSession, len(sm.sessionsMTPM))
	for k, v := range sm.sessionsMTPM {
		copyMap[k] = v
	}

	return copyMap
}

func (sm *SessionManager) GetAllMTPMProgress() map[string]int64 {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	copyMap := make(map[string]int64, len(sm.sessionsMTPM))
	for k, v := range sm.sessionsMTPM {
		copyMap[k] = v.GetSessionProgress()
	}
	return copyMap
}

// DeleteSSE removes a session manually.
func (sm *SessionManager) DeleteMTPM(token string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessionsMTPM, token)
}

// updateLoop runs every 10 seconds to remove expired sessions and update the on-going simulations.
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
	sm.mu.RLock()
	sessions := make([]*engine.TrackedMTPMSession, 0, len(sm.sessionsMTPM))
	for _, v := range sm.sessionsMTPM {
		sessions = append(sessions, v)
	}
	sm.mu.RUnlock()

	for _, v := range sessions {
		v.UpdateAllSubscribers()
	}
}
