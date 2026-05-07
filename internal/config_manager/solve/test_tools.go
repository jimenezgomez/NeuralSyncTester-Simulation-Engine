package config_manager

import (
	"encoding/json"
	"fmt"
	"os"

	config_manager "github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/config_manager/types"
)

type Combination struct {
	K []int
	N []int
}

func FilterByLastK(combinations []Combination, x int) []Combination {
	var result []Combination
	for _, c := range combinations {
		if len(c.K) > 0 && c.K[len(c.K)-1] == x {
			result = append(result, c)
		}
	}
	return result
}

// buildK reconstructs the full K array from Kh and the N array.
// K[H-1] = Kh, and K[i] = K[i+1] * N[i+1] for i < H-1.
func buildK(N []int, Kh int) []int {
	H := len(N)
	K := make([]int, H)
	K[H-1] = Kh
	for i := H - 2; i >= 0; i-- {
		K[i] = K[i+1] * N[i+1]
	}
	return K
}

// buildNFullOverlap reconstructs the full N array from N0 and the K array.
// N[i] = K[i-1] for i < H-1.
func buildNFullOverlap(K []int, N0 int) []int {
	H := len(K)
	N := make([]int, H)
	N[0] = N0
	for i := 1; i < H; i++ {
		N[i] = K[i-1]
	}
	return N
}

// buildNPartialOverlap reconstructs the full N array from N0 and the K array.
// N[i] = K[i-1] - K[i] + 1 for i < H-1.
func buildNPartialOverlap(K []int, N0 int) []int {
	H := len(K)
	N := make([]int, H)
	N[0] = N0
	for i := 1; i < H; i++ {
		N[i] = K[i-1] - K[i] + 1
	}
	return N
}

// GenerateCombinationsConfigFile creates a JSON batch settings file from a list of K/N combinations.
// All fields not derivable from the combinations must be provided explicitly as parameters.
func GenerateCombinationsConfigFile(
	filename string,
	combinations []Combination,
	base config_manager.BaseBatchSettings,
) error {
	// Build the k_n_combinations field: each entry is [[k0,k1,...], [n0,n1,...]].
	knCombinations := make([][][]int, len(combinations))
	for i, c := range combinations {
		knCombinations[i] = [][]int{c.K, c.N}
	}

	raw := config_manager.BatchSettingsFileWithCombinations{
		BaseBatchSettings: base,
		KNCombinations:    knCombinations,
	}

	data, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file %q: %w", filename, err)
	}

	return nil
}
