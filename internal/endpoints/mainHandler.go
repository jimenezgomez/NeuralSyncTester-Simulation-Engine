package endpoints

import (
	"log"
	"net/http"
	"time"

	config_manager "github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/config_manager/load"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/dbmanager"
	graphendpoints "github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/endpoints/graph_endpoints"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine/sync"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/session_manager"
)

var sessionManager = session_manager.NewSessionManager(timeToLive)
var timeToLive = 1 * time.Minute

func RunServerMode(simConfig config_manager.SimulationConfig, trackingConfig config_manager.TrackingConfig, queryManager *dbmanager.QueryManager) {
	settings := engine.NewMTPMSettings([]int{3}, []int{3}, 3, 1, 1, "HEBBIAN", "NO_OVERLAP")
	trackedState := engine.NewTrackedSession(settings, 1_000_000, 256)
	sessionManager.AddMTPM(trackedState.UID, trackedState)
	// fmt.Println(trackedState.UID)
	go func() {
		for i := 0; i < 1_000_000; i++ {
			sync.SimulateTrackedSync(trackedState, simConfig, trackingConfig)
			trackedState.AddProgress()
		}
	}()

	// Register endpoints
	http.HandleFunc("/new_sse", NewSSESessionHandler)
	http.HandleFunc("/sse", SSEHandler)
	http.HandleFunc("/all-sessions", HandleMTPMSessions(sessionManager))
	http.HandleFunc("/all-sessions/progress", HandleMTPMSessionsProgress(sessionManager))
	http.Handle("/api/graph3d", graphendpoints.Graph3DHandler(queryManager))

	addr := ":8080"
	log.Printf("Server listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
