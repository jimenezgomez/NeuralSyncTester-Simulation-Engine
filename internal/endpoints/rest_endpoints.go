package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/session_manager"
)

func HandleMTPMSessions(sm *session_manager.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		snapshot := sm.GetAllMTPMSessions()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(snapshot)
	}
}
func HandleMTPMSessionsProgress(sm *session_manager.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		snapshot := sm.GetAllMTPMProgress()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(snapshot)
	}
}
