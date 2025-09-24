package attacks

import (
	"math"

	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/pkg/tpm/tpm_core"
)

func NewGeomAttack(attackSettings AttackSettings, simulationInstance engine.SimulationInstance) AttackInstance {
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
		attackerExec:       ExecGeomAttack,
		attackerCheck:      CheckGeomAttack,
		attackerStates:     attackers,
	}
	return attackInstance
}

func ExecGeomAttack(settings AttackSettings, sessionState AttackInstance, output_A, output_B int, input_stimulus [][]int) {
	for _, v := range sessionState.attackerStates {
		v.Stimulate(settings.MTPMSettings, input_stimulus)
		learnGeomAttackReduced(settings, v, output_A, output_B)
	}
}

func CheckGeomAttack(settings AttackSettings, sessionState AttackInstance) int {
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

func learnGeomAttackReduced(settings AttackSettings, attacker *AttackerState, output_A, output_B int) {

	if attacker.NetworkOutput != output_A {
		flipLowestLocalField(settings.MTPMSettings, attacker)
	}
	attacker.Learn(settings.MTPMSettings, output_A) //"... then the attacker updates C by the usual learning rule."

}

func flipLowestLocalField(settings engine.MTPMSettings, sessionState *AttackerState) {
	lastLayerIndex := settings.H - 1
	lastLayer := sessionState.Weights[lastLayerIndex]

	minFieldIndexes := []int{}
	minAbsField := math.MaxFloat64
	for neuronIndex := 0; neuronIndex < settings.K[lastLayerIndex]; neuronIndex++ {
		absRawLocalField := math.Abs(
			tpm_core.NeuronLocalFieldRaw(settings.N[lastLayerIndex], lastLayer[neuronIndex],
				sessionState.MTPMState.InputBuffer[lastLayerIndex][neuronIndex]))
		if absRawLocalField < minAbsField {
			minFieldIndexes = []int{neuronIndex}
			minAbsField = absRawLocalField
			continue
		}
		if absRawLocalField == minAbsField {
			minFieldIndexes = append(minFieldIndexes, neuronIndex)
		}
	}

	flipIndex := minFieldIndexes[tpm_core.CryptoRandIntn(len(minFieldIndexes))]

	sessionState.OutputBuffer[lastLayerIndex][flipIndex] *= -1
}
