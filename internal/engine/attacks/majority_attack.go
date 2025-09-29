package attacks

import (
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine"
)

//Original publication:
//L. N. Shacham, E. Klein, R. Mislovaty, I. Kanter, and W. Kinzel. Cooperating attackers in neural cryptography.
// Phys. Rev. E, 69(6):066137, 2004.

func NewMajorityAttack(attackSettings AttackSettings, simulationInstance engine.SimulationInstance) AttackInstance {
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
		attackerExec:       ExecMajorityAttack,
		attackerCheck:      CheckMajorityAttack,
		attackerStates:     attackers,
	}
	return attackInstance
}

// "Additionally, E replaces the majority attack by the geometric attack
// in the first 100 time steps of the synchronization process."
// A. Ruttor, “Neural synchronization and cryptography,” Jan. 2006.
// Note: The original attack proposal says "after a waiting time of about 1/3 of the entire synhronization time"
var START_PHASE_THRESHOLD = 100

func ExecMajorityAttack(settings *AttackSettings, sessionState *AttackInstance, output_A, output_B int, input_stimulus [][]int) {
	operationIndex := sessionState.StimulateIterations & 1

	//"In every odd time step we perform the regular skipping attack,
	//and in every even time step we perform a majority-flipping procedure (...)"
	//L. N. Shacham, E. Klein, R. Mislovaty, I. Kanter, and W. Kinzel. Cooperating attackers in neural cryptography.
	// Phys. Rev. E, 69(6):066137, 2004.
	if sessionState.StimulateIterations < START_PHASE_THRESHOLD || operationIndex == 1 {
		for _, attacker := range sessionState.attackerStates {
			attacker.Stimulate(settings.MTPMSettings, input_stimulus)
			learnGeomAttack(settings, attacker, output_A, output_B)
		}
	} else {
		referenceCombination := getMostCommonCombination(*settings, sessionState.attackerStates, input_stimulus, output_A)
		for _, attacker := range sessionState.attackerStates {
			attacker.Stimulate(settings.MTPMSettings, input_stimulus)
			attacker.LearnWithFullReference(settings.MTPMSettings, output_A, output_B, referenceCombination)
		}
	}

}

func CheckMajorityAttack(settings *AttackSettings, sessionState *AttackInstance) int {
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

type outputCombinationCount struct {
	Combination *[][]int
	Count       int
}

func getMostCommonCombination(tpmSettings AttackSettings, tpmStates []*AttackerState, input_stimulus [][]int, output_A int) *[][]int {
	counts := make(map[uint64]*outputCombinationCount)
	lastLayerIndex := tpmSettings.H - 1
	var maxKey uint64
	maxCount := 0
	for _, tpm := range tpmStates {
		tpm.Stimulate(tpmSettings.MTPMSettings, input_stimulus)
		if tpm.NetworkOutput != output_A {
			flipLowestLocalField(tpmSettings.MTPMSettings, tpm) //This way all attackers have the same output as Alice
		}

		key := encodeOutputsAsBits(tpm.OutputBuffer[lastLayerIndex])
		if entry, ok := counts[key]; ok {
			entry.Count++

			if counts[key].Count > maxCount {
				maxCount = counts[key].Count
				maxKey = key
			}

		} else {
			counts[key] = &outputCombinationCount{
				Combination: &tpm.OutputBuffer, // original list reference
				Count:       1,
			}
			if maxCount == 0 {
				maxCount = 1
				maxKey = key
			}
		}
	}

	if counts[maxKey] == nil {
		panic("maxKey was not set correctly")
	}
	return counts[maxKey].Combination
}

func encodeOutputsAsBits(list []int) uint64 {
	var bits uint64
	for _, v := range list {
		bits <<= 1
		if v == 1 {
			bits |= 1
		}
	}
	return bits
}
