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
	AttackerScores      interface{}
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
	scoresBytes, err := json.Marshal(result.AttackerScores)
	if err != nil {
		return AttackSessionLog{}, err
	}

	return AttackSessionLog{
		NetworkSize:         engine.GetDataSize(mtpmSettings),
		FirstK:              firstOrZero(mtpmSettings.K),
		FirstN:              firstOrZero(mtpmSettings.N),
		LastK:               lastOrZero(mtpmSettings.K),
		LastN:               lastOrZero(mtpmSettings.N),
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
		AttackerScores:      scoresBytes,
		SessionStatus:       result.SessionStatus,
	}, nil
}
