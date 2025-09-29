package attacks

import (
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine"
)

func NewSimpleAttack(attackSettings AttackSettings, simulationInstance engine.SimulationInstance) AttackInstance {
	attackerCount := attackSettings.attackerLimit
	attackers := make([]*AttackerState, attackerCount)
	for i := 0; i < attackerCount; i++ {
		attackers[i] = &AttackerState{
			MTPMState:     engine.NewMTPMState(attackSettings.MTPMSettings),
			attackerScore: 0, weightScore: 0}
	}
	attackInstance := AttackInstance{
		SimulationInstance: simulationInstance,
		attackerCount:      attackerCount,
		attackerExec:       ExecSimpleAttack,
		attackerCheck:      CheckSimpleAttack,
		attackerStates:     attackers,
	}
	return attackInstance
}

func ExecSimpleAttack(settings *AttackSettings, sessionState *AttackInstance, output_A, output_B int, input_stimulus [][]int) {
	for _, v := range sessionState.attackerStates {
		v.Stimulate(settings.MTPMSettings, input_stimulus)
		v.LearnWithOutputs(settings.MTPMSettings, output_A, output_B)
	}
}

func CheckSimpleAttack(settings *AttackSettings, sessionState *AttackInstance) int {
	state := 0
	if engine.CompareWeights(settings.MTPMSettings, sessionState.StateA, sessionState.StateB) {
		state = 1
	}
	for _, attacker := range sessionState.attackerStates {
		if engine.CompareWeights(settings.MTPMSettings, sessionState.StateA, attacker.MTPMState) {
			state = -1
			break
		}
	}

	return state
}
