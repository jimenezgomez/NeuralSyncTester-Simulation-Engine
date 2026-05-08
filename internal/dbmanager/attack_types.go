package dbmanager

import (
	"encoding/json"
	"time"

	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine/attacks"
)

type AttackSessionLog struct {
	NetworkSize         int
	FirstK              int
	FirstN              int
	LastK               int
	LastN               int
	TotalK              int
	TotalN              int
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
	AttackType          string
	AttackerCountLimit  int
	AttackerOverlaps    interface{}
	BestAttackerOverlap float64
	BestAttackerScore   float64
	FinalState          interface{}
	SessionStatus       string
}

func NewAttackSessionLog(result attacks.AttackResult) (AttackSessionLog, error) {
	attSettings := result.Settings
	mtpmSettings := attSettings.MTPMSettings

	// marshal K and N (safe for JSONB)
	kBytes, err := json.Marshal(mtpmSettings.K)
	if err != nil {
		return AttackSessionLog{}, err
	}
	nBytes, err := json.Marshal(mtpmSettings.N)
	if err != nil {
		return AttackSessionLog{}, err
	}
	scoresBytes, err := json.Marshal(result.AttackerOverlaps)
	if err != nil {
		return AttackSessionLog{}, err
	}

	return AttackSessionLog{
		NetworkSize:         engine.GetDataSize(mtpmSettings),
		FirstK:              firstOrZero(mtpmSettings.K),
		FirstN:              firstOrZero(mtpmSettings.N),
		LastK:               lastOrZero(mtpmSettings.K),
		LastN:               lastOrZero(mtpmSettings.N),
		TotalK:              sumArray(mtpmSettings.K),
		TotalN:              sumArray(mtpmSettings.N),
		StartTime:           result.StartTime,
		EndTime:             result.EndTime,
		StimulateIterations: result.FinalState.StimulateIterations,
		LearnIterations:     result.FinalState.LearnIterations,
		K:                   kBytes,
		N:                   nBytes,
		L:                   mtpmSettings.L,
		M:                   mtpmSettings.M,
		H:                   mtpmSettings.H,
		LearnRule:           mtpmSettings.LearnRule,
		Scenario:            mtpmSettings.Scenario,
		AttackType:          attSettings.AttackType,
		AttackerCountLimit:  attSettings.AttackerLimit,
		AttackerOverlaps:    scoresBytes,
		BestAttackerOverlap: result.BestAttackerOverlap,
		BestAttackerScore:   result.BestAttackerScore,
		SessionStatus:       result.SessionStatus,
	}, nil
}

func sumArray(nums []int) int {
	total := 0
	// _ ignores the index, num is the current value
	for _, num := range nums {
		total += num
	}
	return total
}
