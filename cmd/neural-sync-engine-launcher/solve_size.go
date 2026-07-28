package cmd

import (
	"fmt"
	"path/filepath"

	config_solver "github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/config_manager/solve"
	config_manager "github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/config_manager/types"
	"github.com/spf13/cobra"
)

var sizeSolverCmd = &cobra.Command{
	Use:   "solve",
	Short: "(EXPERIMENTAL) Run the NeuralSyncTester config solver auxiliary tool",
	Long: `Creates valid configurations for a given size to use in NeuralSyncTester engine. 
	Experimental. To use, you must recompile with you own settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		// solverOutput := config_solver.SolveSize_NoOverlap(1, 3, 0, 3)
		targetSize := 3 * 1000 // Using ruttor's 2006 thesis values in chapter 3-4 for K and N respectively
		maxLayercount := 5

		minN0 := 2
		maxN0 := targetSize / 2

		var fullOverlap_solverOutput []config_solver.Combination
		var partialOverlap_solverOutput []config_solver.Combination
		var noOverlap_solverOutput []config_solver.Combination

		for i := 2; i <= maxLayercount; i++ {
			// solverOutput := config_solver.SolveSize_Overlap(2, targetSize, 2, targetSize/2, "PARTIAL OVERLAP")
			solverOutput := config_solver.SolveSize_Overlap(i, targetSize, minN0, maxN0, "FULL OVERLAP")
			fullOverlap_solverOutput = append(fullOverlap_solverOutput, solverOutput...)

			solverOutput = config_solver.SolveSize_Overlap(i, targetSize, minN0, maxN0, "PARTIAL OVERLAP")
			partialOverlap_solverOutput = append(partialOverlap_solverOutput, solverOutput...)

			solverOutput = config_solver.SolveSize_NoOverlap(i, targetSize, minN0, maxN0) //This function uses min/max Kh instead, but same values work
			noOverlap_solverOutput = append(noOverlap_solverOutput, solverOutput...)
		}
		baseConfig := config_manager.BaseBatchSettings{
			Scenario:        "full_overlap",
			MaxSessionCount: 10,
			MaxIterations:   100000,
			MaxWorkerCount:  6,
			LearnRules: []string{
				// Only using random walk because of fig 3.24 on Ruttor 2006 thesis
				// "Hebbian",
				// "Anti-Hebbian",
				"Random-Walk",
			},
			// You must specify the type (e.g., []int) for these slices
			MConfigs: []int{1},
			LConfigs: []int{3, 5, 7, 10},
		}

		base_filename := filepath.Join("config", "experimento_slim", "multi_layer_3k_slim")
		total_combi_count := 0

		lastKAmount := 3
		fullOverlap_solverOutput_filtered := config_solver.FilterByLastK(fullOverlap_solverOutput, lastKAmount)
		partialOverlap_solverOutput_filtered := config_solver.FilterByLastK(partialOverlap_solverOutput, lastKAmount)
		noOverlap_solverOutput_filtered := config_solver.FilterByLastK(noOverlap_solverOutput, lastKAmount)

		baseConfig.Scenario = "full_overlap"
		output_filename := base_filename + "/" + baseConfig.Scenario + ".json"
		config_solver.GenerateCombinationsConfigFile(output_filename, fullOverlap_solverOutput_filtered, baseConfig)
		total_combi_count += len(fullOverlap_solverOutput_filtered)

		baseConfig.Scenario = "partial_overlap"
		output_filename = base_filename + "/" + baseConfig.Scenario + ".json"
		config_solver.GenerateCombinationsConfigFile(output_filename, partialOverlap_solverOutput_filtered, baseConfig)
		total_combi_count += len(partialOverlap_solverOutput_filtered)

		baseConfig.Scenario = "no_overlap"
		output_filename = base_filename + "/" + baseConfig.Scenario + ".json"
		config_solver.GenerateCombinationsConfigFile(output_filename, noOverlap_solverOutput_filtered, baseConfig)
		total_combi_count += len(noOverlap_solverOutput_filtered)

		fmt.Printf("Total combination count: %d\n", total_combi_count)
		fmt.Printf("Total combination count with base parameters: %d\n", total_combi_count*baseConfig.MaxSessionCount*len(baseConfig.LearnRules)*len(baseConfig.LConfigs)*len(baseConfig.MConfigs))
	},
}
