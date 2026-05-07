package cmd

import (
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/endpoints"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the NeuralSyncTester server",
	Long:  `Starts the HTTP server for SSE and other endpoints for NeuralSyncTester.`,
	Run: func(cmd *cobra.Command, args []string) {
		endpoints.RunServerMode(GlobalSimulationConfig, GlobalTrackingConfig, queryManager, sessionManager)
	},
}
