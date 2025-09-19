package config_manager

import (
	"fmt"

	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/pkg/tpm/tpm_stimHandlers"
)

type BatchScenarioWrapper struct {
	Scenario string `json:"scenario"`
}

type BaseBatchSettings struct {
	Scenario        string   `json:"scenario"`
	MaxSessionCount int      `json:"max_session_count"`
	MaxIterations   int      `json:"max_iterations"`
	MaxWorkerCount  int      `json:"max_worker_count"`
	LearnRules      []string `json:"learn_rules"`
	MConfigs        []int    `json:"m_configs"`
	LConfigs        []int    `json:"l_configs"`
}

type BatchOverlappedSettings struct {
	BaseBatchSettings
	KConfigs  [][]int `json:"k_configs"`
	N0Configs []int   `json:"n0_configs"`
}

type BatchNonOverlappedSettings struct {
	BaseBatchSettings
	KlastConfigs []int   `json:"klast_configs"`
	NConfigs     [][]int `json:"n_configs"`
}

type BatchSettingsLoader interface {
	UnwrapBatchSettings(configList *[]engine.MTPMSettings, learnRule string, m, l int)
	GetBase() BaseBatchSettings
}

func (settings BatchOverlappedSettings) UnwrapBatchSettings(configList *[]engine.MTPMSettings, learnRule string, m, l int) {
	for _, k := range settings.KConfigs {
		for _, n_0 := range settings.N0Configs {
			var N []int
			switch settings.Scenario {
			case "FULL_OVERLAP":
				N = tpm_stimHandlers.FullOverlapTPM{}.CreateStimulationStructure(k, n_0)
			case "PARTIAL_OVERLAP":
				N = tpm_stimHandlers.PartialOverlapTPM{}.CreateStimulationStructure(k, n_0)
			default:
				fmt.Printf("Error: %s is not an overlapped scenario\n", settings.Scenario)
				return
			}
			h := len(N)
			tpmInstanceSettings := engine.NewMTPMSettings(k, N, l, m, h, learnRule, settings.Scenario)
			*configList = append(*configList, tpmInstanceSettings)
		}
	}
}
func (s BatchOverlappedSettings) GetBase() BaseBatchSettings { return s.BaseBatchSettings }

func (settings BatchNonOverlappedSettings) UnwrapBatchSettings(configList *[]engine.MTPMSettings, learnRule string, m, l int) {
	for _, n := range settings.NConfigs {
		for _, k_last := range settings.KlastConfigs {
			K := tpm_stimHandlers.NoOverlapTPM{}.CreateStimulationStructure(n, k_last)
			h := len(K)
			tpmInstanceSettings := engine.NewMTPMSettings(K, n, l, m, h, learnRule, settings.Scenario)
			*configList = append(*configList, tpmInstanceSettings)
		}
	}
}
func (s BatchNonOverlappedSettings) GetBase() BaseBatchSettings { return s.BaseBatchSettings }
