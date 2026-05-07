package config_manager

import (
	"encoding/json"
	"fmt"
	"os"

	config_manager "github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/config_manager/types"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine"
)

// LoadBatchSettingsFromFile_Combinations loads a batch settings JSON file that uses
// k_n_combinations to define K/N pairs, and expands them into all combinations of
// K/N x M x L x LearnRule, returning a flat []engine.MTPMSettings.
func LoadBatchSettingsFromFile_Combinations(filename string) ([]engine.MTPMSettings, *config_manager.BaseBatchSettings, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file %q: %w", filename, err)
	}

	var raw config_manager.BatchSettingsFileWithCombinations
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, nil, fmt.Errorf("failed to parse JSON in %q: %w", filename, err)
	}

	// Validate k_n_combinations entries — each must be a pair [K, N].
	for i, combo := range raw.KNCombinations {
		if len(combo) != 2 {
			return nil, nil, fmt.Errorf(
				"k_n_combinations[%d]: expected [K, N] pair (2 slices), got %d slice(s)",
				i, len(combo),
			)
		}
	}

	// Expand all combinations: k_n x m x l x learn_rule.
	var settings []engine.MTPMSettings
	for _, combo := range raw.KNCombinations {
		k := combo[0]
		n := combo[1]
		h := len(k)
		for _, m := range raw.MConfigs {
			for _, l := range raw.LConfigs {
				for _, rule := range raw.LearnRules {
					tpmInstanceSettings := engine.NewMTPMSettings(k, n, l, m, h, rule, raw.Scenario)
					settings = append(settings, tpmInstanceSettings)
				}
			}
		}
	}

	return settings, &raw.BaseBatchSettings, nil
}
