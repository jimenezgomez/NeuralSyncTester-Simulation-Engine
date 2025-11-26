package config_manager

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	config_manager "github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/config_manager/types"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine"
)

// LoadBatchSettingsFromFile creates a list of settings from a filename
// Returns: A list of MTPMSettings, The base settings (shared by all elements on the list
func LoadBatchSettingsFromFile(filename string) ([]engine.MTPMSettings, *config_manager.BaseBatchSettings, error) {
	// read file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file: %w", err)
	}

	// first unmarshal only the type
	var wrapper config_manager.BatchScenarioWrapper
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, nil, fmt.Errorf("failed to parse wrapper: %w", err)
	}

	// unmarshal into the correct struct based on type
	switch strings.ToUpper(wrapper.Scenario) {
	case "FULL_OVERLAP", "PARTIAL_OVERLAP":
		var s config_manager.BatchOverlappedSettings
		if err := json.Unmarshal(data, &s); err != nil {
			return nil, nil, fmt.Errorf("failed to parse overlapped settings: %w", err)
		}
		return ListBatchSettings(s), &s.BaseBatchSettings, nil
	case "NO_OVERLAP":
		var s config_manager.BatchNonOverlappedSettings
		if err := json.Unmarshal(data, &s); err != nil {
			return nil, nil, fmt.Errorf("failed to parse non-overlapped settings: %w", err)
		}
		return ListBatchSettings(s), &s.BaseBatchSettings, nil
	default:
		return nil, nil, fmt.Errorf("unknown scenario: %s", wrapper.Scenario)
	}
}

// ListBatchSettings uses the loader provided by LoadBatchSettings to unwrap the settings, depending if its Overlapped or non Overlapped
func ListBatchSettings(loader config_manager.BatchSettingsLoader) []engine.MTPMSettings {
	configList := make([]engine.MTPMSettings, 0)
	baseSettings := loader.GetBase()
	for _, rule := range baseSettings.LearnRules {
		for _, m := range baseSettings.MConfigs {
			for _, l := range baseSettings.LConfigs {
				loader.UnwrapBatchSettings(&configList, rule, m, l)
			}
		}
	}
	return configList
}
