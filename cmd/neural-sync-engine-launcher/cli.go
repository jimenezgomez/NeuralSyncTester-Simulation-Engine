package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/session_manager"
	"github.com/spf13/cobra"
)

var timeToLive = 1 * time.Minute

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Run a simulation in CLI mode (no SSE)",
	Run: func(cmd *cobra.Command, args []string) {
		settings := engine.NewMTPMSettings([]int{3}, []int{3}, 3, 1, 1, "HEBBIAN", "NO_OVERLAP")
		trackedState := engine.NewTrackedState(settings)
		sessionManager := session_manager.NewSessionManager(timeToLive)
		sessionManager.AddMTPM(trackedState.UID, trackedState)
		ch := make(chan []byte, 10) // small buffer, prevent blocking
		trackedState.Subscribe(ch)
		defer trackedState.Unsubscribe(ch)

		go func() {
			for i := 0; i < 1_000_000; i++ {
				engine.SimulateTrackedSync(trackedState)
			}
		}()

		for msg := range ch {
			var snapshot engine.SimulationInstance
			if err := json.Unmarshal(msg, &snapshot); err != nil {
				fmt.Println("failed to unmarshal snapshot:", err)
				continue
			}

			// Now you can access fields directly
			fmt.Printf("%s\n", trackedState.StartTime.Format(time.ANSIC))
			fmt.Printf("%s", snapshot.PrettyPrint())
		}

	},
}
