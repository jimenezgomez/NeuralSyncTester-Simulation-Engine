package attacks

import (
	"sort"
	"strings"
	"time"

	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/pkg/tpm/tpm_core"
)

const ITERATION_LIMIT = 100_000
const TRACKING_SKIP_ITERATIONS = 150 //Skip n iterations before sending a new attackerState
const DEFAULT_ATT_LIMIT = 100        //We will use this many attackers. (All implemented attacks use up to the limit - other attacks may use this as a limit and dynamically spawn new attackers)
const STORE_TOP_ATT_LIMIT = 10       //only save the best 10 attackers

func RunTrackedAttack(trackedState *engine.TrackedMTPMSession, attackType string) AttackResult {
	settings := trackedState.GetSettings()

	attackSettings := AttackSettings{
		MTPMSettings:  settings,
		AttackerLimit: DEFAULT_ATT_LIMIT,
		AttackType:    attackType,
	}

	simulationInstance := engine.SimulationInstance{
		SimulationState: engine.SimulationState{
			StateA: engine.NewMTPMState(settings),
			StateB: engine.NewMTPMState(settings)},
		SimulationProgress: engine.SimulationProgress{
			StimulateIterations: 0,
			LearnIterations:     0,
		},
	}

	attackInstance := CreateAttackInstance(simulationInstance, attackSettings)
	if attackInstance == nil {
		panic("The requested attack type was not found. - " + attackSettings.AttackType)
	}

	startInstance := simulationInstance.DeepCopy() //We store a copy of the initial state to store it later in the database
	trackedState.UpdateSnapshot(startInstance)
	trackedState.StartTime = time.Now()
	syncReached := tpm_core.CompareWeights(settings.H, settings.K, settings.N,
		simulationInstance.StateA.Weights, simulationInstance.StateB.Weights)

	checkResult := 0
	skipCounter := TRACKING_SKIP_ITERATIONS
	for !syncReached {
		if simulationInstance.StimulateIterations > ITERATION_LIMIT {
			break
		}

		input_stimulus := tpm_core.CreateRandomStimulusArray(settings.K[0], settings.N[0], settings.M)
		simulationInstance.StateA.Stimulate(settings, input_stimulus)
		simulationInstance.StateB.Stimulate(settings, input_stimulus)
		// StimulateAllAttackers(attackSettings, attackInstance, input_stimulus)
		// Here we should stimulate the attackers, but we don't really need to unless we need to update their weights
		// If no learning is done, the behaviour of the TPMs should not change! (they keep their weights w/o change)
		// We stimulate all attackers on attackerExec, before actually performing the attack
		simulationInstance.StimulateIterations += 1

		if simulationInstance.StateA.NetworkOutput == simulationInstance.StateB.NetworkOutput {
			attackInstance.attackerExec(
				&attackSettings, attackInstance,
				attackInstance.StateA.NetworkOutput, attackInstance.StateB.NetworkOutput,
				input_stimulus)
			simulationInstance.StateA.Learn(settings, simulationInstance.StateB.NetworkOutput)
			simulationInstance.StateB.Learn(settings, simulationInstance.StateA.NetworkOutput)
			simulationInstance.LearnIterations += 1
		}

		checkResult = attackInstance.attackerCheck(&attackSettings, attackInstance)
		syncReached = checkResult != 0

		//TRACKING
		if trackedState.GetSubCount() > 0 {
			// Only bother tracking/publishing if someone is listening
			if skipCounter == 0 {
				snapshot := simulationInstance.DeepCopy()
				trackedState.UpdateSnapshot(snapshot)
				skipCounter = TRACKING_SKIP_ITERATIONS
			}
			skipCounter--
		}

	}

	status := "LIMIT_REACHED"
	if checkResult > 0 {
		status = "ON_SYNC"
	}
	if checkResult < 0 {
		status = "ATTACK_SUCCESS"
	}

	//Get the best attackers, they will be saved in the db
	SetAttackerWeightScores(attackSettings, *attackInstance)
	topAttackers := GetTopAttackersByWeight(*attackInstance)
	topScores := make([]float64, STORE_TOP_ATT_LIMIT)
	for i := range STORE_TOP_ATT_LIMIT {
		topScores[i] = topAttackers[i].weightScore
	}

	result := AttackResult{
		Settings:       attackSettings,
		FinalState:     simulationInstance,
		SessionStatus:  status,
		StartTime:      trackedState.StartTime,
		EndTime:        time.Now(),
		AttackerScores: topScores,
	}

	return result
}

// Helper function - used for tests only
func StimulateAllAttackers(settings AttackSettings, instance AttackInstance, input_stimulus [][]int) {
	for _, v := range instance.attackerStates {
		v.Stimulate(settings.MTPMSettings, input_stimulus)
	}
}

func SetAttackerWeightScores(settings AttackSettings, instance AttackInstance) {
	for _, attacker := range instance.attackerStates {
		attacker.weightScore = tpm_core.CosineSimWeights(settings.H, settings.K, settings.N, instance.StateA.Weights, attacker.Weights)
	}
}

func GetTopAttackersByWeight(instance AttackInstance) []*AttackerState {
	// Make a copy of the slice of pointers
	attackers := make([]*AttackerState, len(instance.attackerStates))
	copy(attackers, instance.attackerStates)

	// Sort in-place by weightScore (highest first)
	sort.Slice(attackers, func(i, j int) bool {
		return attackers[i].weightScore > attackers[j].weightScore
	})

	return attackers
}

func CreateAttackInstance(simulationInstance engine.SimulationInstance, attackSettings AttackSettings) *AttackInstance {
	switch strings.ToUpper(attackSettings.AttackType) {
	case "SIMPLE", "NAIVE":
		attackInstance := NewSimpleAttack(attackSettings, simulationInstance)
		return &attackInstance
	case "GEOMETRIC":
		attackInstance := NewGeomAttack(attackSettings, simulationInstance)
		return &attackInstance
	case "MAJORITY", "MAJORITY-FLIPPING", "MAJORITY FLIPPING":
		attackInstance := NewMajorityAttack(attackSettings, simulationInstance)
		return &attackInstance
	}
	return nil
}
