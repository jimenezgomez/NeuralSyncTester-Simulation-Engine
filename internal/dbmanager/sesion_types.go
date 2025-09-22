package dbmanager

import (
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

func NewLogFromResult(res engine.SimulationResult) SyncSessionLog {
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
		K:                   res.Settings.K, // full object goes to JSONB
		N:                   res.Settings.N,
		L:                   res.Settings.L,
		M:                   res.Settings.M,
		H:                   res.Settings.H,
		LearnRule:           res.Settings.LearnRule,
		Scenario:            res.Settings.Scenario,
		InitialState:        res.InitialState,
		FinalState:          res.FinalState,
		SessionStatus:       res.SessionStatus,
	}
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
