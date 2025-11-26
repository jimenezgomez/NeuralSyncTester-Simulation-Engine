package attacks

import (
	"time"

	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine"
)

type AttackSettings struct {
	engine.MTPMSettings
	AttackerLimit int
	AttackType    string
}

type AttackerState struct {
	engine.MTPMState
	attackerScore int     //The score based on the attack type
	weightScore   float64 //The score based on the dot product

}

type AttackInstance struct {
	engine.SimulationInstance
	attackerCount  int
	attackerStates []*AttackerState
	attackerExec   AttackExec
	attackerCheck  AttackCheck
}

type AttackResult struct {
	Settings       AttackSettings
	FinalState     engine.SimulationInstance
	AttackerScores []float64
	SessionStatus  string
	StartTime      time.Time
	EndTime        time.Time
}

type AttackExec func(settings *AttackSettings, sessionState *AttackInstance, output_A, output_B int, input_stimulus [][]int)

type AttackCheck func(settings *AttackSettings, sessionState *AttackInstance) int
