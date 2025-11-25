package dbmanager

import (
	"encoding/json"
	"time"

	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine"
)

type SyncSessionLog struct {
	NetworkSize         int
	FirstK              int
	FirstN              int
	LastK               int
	LastN               int
	StartTime           time.Time
	EndTime             time.Time
	StimulateIterations int
	LearnIterations     int
	K                   interface{} // will be marshaled into JSONB
	N                   interface{}
	L                   int
	M                   int
	H                   int
	LearnRule           string
	Scenario            string
	InitialState        interface{}
	FinalState          interface{}
	SessionStatus       string
}

func NewLogFromResult(res engine.SimulationResult) (SyncSessionLog, error) {
	// marshal K and N (safe for JSONB)
	mtpmSettings := res.Settings

	kBytes, err := json.Marshal(mtpmSettings.K)
	if err != nil {
		return SyncSessionLog{}, err
	}
	nBytes, err := json.Marshal(mtpmSettings.N)
	if err != nil {
		return SyncSessionLog{}, err
	}

	initialState, err := json.Marshal(res.InitialState)
	if err != nil {
		return SyncSessionLog{}, err
	}
	finalState, err := json.Marshal(res.FinalState)
	if err != nil {
		return SyncSessionLog{}, err
	}

	return SyncSessionLog{
		NetworkSize:         engine.GetDataSize(res.Settings),
		FirstK:              firstOrZero(res.Settings.K),
		FirstN:              firstOrZero(res.Settings.N),
		LastK:               lastOrZero(res.Settings.K),
		LastN:               lastOrZero(res.Settings.N),
		StartTime:           res.StartTime,
		EndTime:             res.EndTime,
		StimulateIterations: res.Iterations.StimulateIterations,
		LearnIterations:     res.Iterations.LearnIterations,
		K:                   kBytes,
		N:                   nBytes,
		L:                   res.Settings.L,
		M:                   res.Settings.M,
		H:                   res.Settings.H,
		LearnRule:           res.Settings.LearnRule,
		Scenario:            res.Settings.Scenario,
		InitialState:        initialState,
		FinalState:          finalState,
		SessionStatus:       res.SessionStatus,
	}, err
}

// helpers for extracting first/last from slices
func firstOrZero(arr []int) int {
	if len(arr) > 0 {
		return arr[0]
	}
	return 0
}

func lastOrZero(arr []int) int {
	if len(arr) > 0 {
		return arr[len(arr)-1]
	}
	return 0
}
