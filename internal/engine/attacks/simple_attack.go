package attacks

import (
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine"
)

// In order to reduce fluctuations
// in our simulations we define a successful attacker as one
// which has 0.98 fraction of correct values for the weights at
// the synchronization time between the parties.

func NewSimpleAttack(attackSettings AttackSettings, simulationInstance *engine.SimulationInstance) AttackInstance {
	attackerCount := attackSettings.AttackerLimit
	attackers := make([]*AttackerState, attackerCount)
	for i := 0; i < attackerCount; i++ {
		attackers[i] = &AttackerState{
			MTPMState:     engine.NewMTPMState(attackSettings.MTPMSettings),
			attackerScore: 0, weightOverlap: 0}
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
	for _, attacker := range sessionState.attackerStates {
		attacker.Stimulate(settings.MTPMSettings, input_stimulus)
		attacker.LearnWithOutputs(settings.MTPMSettings, output_A, output_B)
	}
}

func CheckSimpleAttack(settings *AttackSettings, sessionState *AttackInstance) int {
	state := 0
	if engine.CompareWeights(settings.MTPMSettings, sessionState.StateA, sessionState.StateB) {
		state = 1
		for _, attacker := range sessionState.attackerStates {
			if engine.CompareWeights(settings.MTPMSettings, sessionState.StateA, attacker.MTPMState) {
				state = -1
				break
			}
			if attacker.attackerScore >= ATTACKER_SCORE_THRESHOLD {
				state = -2
			}
		}
	}

	return state
}
