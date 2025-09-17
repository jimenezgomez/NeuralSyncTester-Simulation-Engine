package endpoints

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/session_manager"
)

var sessionManager = session_manager.NewSessionManager(timeToLive)
var timeToLive = 1 * time.Minute

func RunServerMode() {
	settings := engine.NewMTPMSettings([]int{3}, []int{3}, 3, 1, 1, "HEBBIAN", "NO_OVERLAP")
	trackedState := engine.NewTrackedState(settings)
	sessionManager.AddMTPM(trackedState.UID, trackedState)
	fmt.Println(trackedState.UID)
	go func() {
		for i := 0; i < 1_000_000; i++ {
			engine.SimulateTrackedSync(trackedState)
		}
	}()

	// Register endpoints
	http.HandleFunc("/new_sse", NewSSESessionHandler)
	http.HandleFunc("/sse", SSEHandler)

	addr := ":8080"
	log.Printf("Server listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
