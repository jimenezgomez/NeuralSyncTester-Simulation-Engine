package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/session_manager"
)

// NewSSESessionHandler handles requests to create a new SSE session
func NewSSESessionHandler(w http.ResponseWriter, r *http.Request) {
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

func SSEHandler(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("uid")
	if uid == "" {
		http.Error(w, "missing 'uid' parameter", http.StatusBadRequest)
		return
	}

	// lookup
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

	// setup SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	// create channel for this subscriber
	ch := make(chan []byte, 10) // small buffer, prevent blocking
	state.Subscribe(ch)
	defer state.Unsubscribe(ch)
	defer sessionManager.DeleteSSE(uid)

	// stream loop
	for {
		select {
		case <-r.Context().Done(): // client disconnected
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		}
	}
}
