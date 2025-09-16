package engine

import (
	"sync"
	"sync/atomic"

	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/pkg/tpm/tpm_core"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/pkg/tpm/tpm_learnRules"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/pkg/tpm/tpm_stimHandlers"
)

type MTPMSettings struct {
	K                   []int
	N                   []int
	L                   int
	M                   int
	H                   int
	LearnRule           string
	Scenario            string
	stimulationHandlers tpm_stimHandlers.TPMStimulationHandlers
	learnRuleHandler    tpm_learnRules.TPMLearnRuleHandler
}

func NewMTPMSettings(K, N []int, L, M, H int, learnRule, scenario string) MTPMSettings {
	var stimHandler tpm_stimHandlers.TPMStimulationHandlers
	var ruleHandler tpm_learnRules.TPMLearnRuleHandler

	switch scenario {
	case "NO_OVERLAP", "NO OVERLAP":
		stimHandler = tpm_stimHandlers.NoOverlapTPM{}
	case "PARTIAL_OVERLAP", "PARTIAL OVERLAP", "PARTIALLY_CONNECTED", "PARTIALLY CONNECTED":
		stimHandler = tpm_stimHandlers.PartialOverlapTPM{}
	case "FULL_OVERLAP", "FULL OVERLAP", "FULLY_CONNECTED", "FULLY CONNECTED":
		stimHandler = tpm_stimHandlers.FullOverlapTPM{}

	}

	switch learnRule {
	case "HEBBIAN":
		ruleHandler = tpm_learnRules.HebbianLearnRule{}
	case "ANTI-HEBBIAN", "ANTIHEBBIAN", "ANTI HEBBIAN":
		ruleHandler = tpm_learnRules.AntiHebbianLearnRule{}
	case "RANDOM-WALK", "RANDOM", "RANDOMWALK", "RANDOM WALK":
		ruleHandler = tpm_learnRules.RandomWalkLearnRule{}
	}

	return MTPMSettings{
		H:                   H,
		K:                   K,
		N:                   N,
		L:                   L,
		M:                   M,
		LearnRule:           learnRule,
		Scenario:            scenario,
		stimulationHandlers: stimHandler,
		learnRuleHandler:    ruleHandler,
	}
}

// Core network state, reusable across different types
type MTPMState struct {
	Weights       [][][]int
	InputBuffer   [][][]int
	OutputBuffer  [][]int
	NetworkOutput int
}

func NewMTPMState(settings MTPMSettings) MTPMState {
	newWeights := make([][][]int, settings.H)
	for layer := 0; layer < settings.H; layer++ {
		newWeights[layer] = tpm_core.CreateRandomLayerWeightsArray(settings.K[layer], settings.N[layer], settings.L)
	}
	return MTPMState{
		Weights:       newWeights,
		InputBuffer:   make([][][]int, settings.H),
		OutputBuffer:  make([][]int, settings.H),
		NetworkOutput: 0,
	}
}

// Simulation-specific extensions
type MTPMStateWithHistory struct {
	MTPMState
	OutputHistory []int
	HistorySize   int
}

type SimulationInstance struct {
	StateA              MTPMState
	StateB              MTPMState
	StimulateIterations int
	LearnIterations     int
}

type TrackedMTPMState struct {
	UID      string
	settings MTPMSettings
	latest   atomic.Value // stores []byte (marshaled snapshot)
	subs     map[chan []byte]struct{}
	subCount int
	subsLock sync.RWMutex
}

func NewTrackedState(settings MTPMSettings) *TrackedMTPMState {
	ts := &TrackedMTPMState{
		UID:      "0",
		settings: settings,
		subs:     make(map[chan []byte]struct{}),
		subCount: 0,
	}
	return ts
}
