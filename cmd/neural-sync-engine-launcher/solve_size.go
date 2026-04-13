package cmd

import (
	"fmt"

	config_solver "github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/config_manager/solve"
	"github.com/spf13/cobra"
)

var sizeSolverCmd = &cobra.Command{
	Use:   "solve",
	Short: "(EXPERIMENTAL) Run the NeuralSyncTester config solver auxiliary tool",
	Long: `Creates valid configurations for a given size to use in NeuralSyncTester engine. 
	Experimental. To use, you must recompile with you own settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		// solverOutput := config_solver.SolveSize_NoOverlap(1, 3, 0, 3)
		solverOutput := config_solver.SolveSize_Overlap(2, 9, 1, 3, "PARTIAL OVERLAP")
		fmt.Println(solverOutput)
	},
}
