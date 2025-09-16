package main

import (
	"fmt"
	"time"

	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine"
	"github.com/spf13/cobra"
)

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Run a simulation in CLI mode (no SSE)",
	Run: func(cmd *cobra.Command, args []string) {
		settings := engine.NewMTPMSettings([]int{3}, []int{3}, 3, 1, 1, "HEBBIAN", "NO_OVERLAP")
		trackedState := engine.NewTrackedState(settings)
		trackedState.Subscribe()
		defer trackedState.Unsubscribe()

		go func() {
			for i := 0; i < 1_000_000; i++ {
				engine.SimulateTrackedSync(trackedState)
			}
		}()

		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			snapshot := trackedState.GetSnapshot()
			fmt.Printf("%s", snapshot.PrettyPrint())
		}
	},
}
