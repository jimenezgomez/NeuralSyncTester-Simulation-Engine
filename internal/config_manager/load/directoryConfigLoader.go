package config_manager

import (
	"fmt"
	"os"
	"path/filepath"

	config_manager "github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/config_manager/types"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine"
)

// Wrapper result allows multiple files
type BatchSettingsResult struct {
	Path              string
	SettingsList      []engine.MTPMSettings
	BaseBatchSettings config_manager.BaseBatchSettings
	Err               error
}

// ScanAndLoadBatchSettings walks through rootDir and all subdirectories,
// loads all JSON files via LoadBatchSettingsFromFile.
func ScanAndLoadBatchSettings(rootDir string) ([]BatchSettingsResult, error) {
	var results []BatchSettingsResult

	err := filepath.WalkDir(rootDir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".json" {
			settingsList, baseBatchSettings, err := LoadBatchSettingsFromFile(path)
			if err != nil {
				return err
			}
			results = append(results, BatchSettingsResult{
				Path:              path,
				SettingsList:      settingsList,
				BaseBatchSettings: *baseBatchSettings,
				Err:               err,
			})
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}

	return results, nil
}

func ScanAndLoadBatchSettings_Combo(rootDir string) ([]BatchSettingsResult, error) {
	var results []BatchSettingsResult

	err := filepath.WalkDir(rootDir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".json" {
			settingsList, baseBatchSettings, err := LoadBatchSettingsFromFile_Combinations(path)
			if err != nil {
				return err
			}
			results = append(results, BatchSettingsResult{
				Path:              path,
				SettingsList:      settingsList,
				BaseBatchSettings: *baseBatchSettings,
				Err:               err,
			})
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}

	return results, nil
}
