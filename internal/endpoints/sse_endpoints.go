package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/session_manager"
)

// NewSSESessionHandler handles requests to create a new SSE session
func NewSSESessionHandler(w http.ResponseWriter, r *http.Request, sessionManager *session_manager.SessionManager) {
	// parse required params
	syncSessionUID := r.URL.Query().Get("session_uid")
	if syncSessionUID == "" {
		http.Error(w, "missing session uid", http.StatusBadRequest)
		return
	}

	_, ok := sessionManager.GetMTPM(syncSessionUID)
	if !ok {
		http.Error(w, "sync session not found", http.StatusNotFound)
		return
	}

	// generate secure token
	token, err := generateToken(16) // 32 hex chars
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	// get client IP (very basic, you can extend this)
	ip := r.RemoteAddr

	// create new session
	session := session_manager.SSESession{
		UID:            token,
		IP:             ip,
		SyncSessionUID: syncSessionUID,
		Expiry:         time.Now().Add(5 * time.Minute), // example expiry
	}

	// add it to the session manager
	sessionManager.AddSSE(token, session)

	// return the token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func SSEHandler(w http.ResponseWriter, r *http.Request, sessionManager *session_manager.SessionManager) {
	uid := r.URL.Query().Get("uid")
	if uid == "" {
		http.Error(w, "missing 'uid' parameter", http.StatusBadRequest)
		return
	}

	session, ok := sessionManager.GetSSE(uid)
	if !ok {
		http.Error(w, "session not found.", http.StatusBadRequest)
		return
	}

	state, ok := sessionManager.GetMTPM(session.SyncSessionUID)
	if !ok {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	ch := make(chan []byte, 10)
	state.Subscribe(ch)
	defer state.Unsubscribe(ch)
	defer sessionManager.DeleteSSE(uid)

	for {
		select {
		case <-r.Context().Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			compressed, err := compressPayload(msg)
			if err != nil {
				continue
			}
			fmt.Fprintf(w, "data: %s\n\n", compressed)
			flusher.Flush()
		}
	}
}
