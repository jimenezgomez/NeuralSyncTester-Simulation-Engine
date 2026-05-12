package attacks

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	config_manager "github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/config_manager/load"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/pkg/tpm/tpm_core"
)

const ATTACKER_SCORE_THRESHOLD float64 = 0.98 //As described in 10.1103/PhysRevE.66.066102:

func RunTrackedAttack(trackedState *engine.TrackedMTPMSession, attackType string, simConfig config_manager.SimulationConfig, trackConfig config_manager.TrackingConfig) AttackResult {
	settings := trackedState.GetSettings()

	attackSettings := AttackSettings{
		MTPMSettings:  settings,
		AttackerLimit: simConfig.DefaultAttackerLimit,
		AttackType:    attackType,
	}

	simulationInstance := engine.SimulationInstance{
		SimulationNumber: int(trackedState.GetSessionProgress()),
		SimulationState: engine.SimulationState{
			StateA: engine.NewMTPMState(settings),
			StateB: engine.NewMTPMState(settings)},
		SimulationProgress: engine.SimulationProgress{
			StimulateIterations: 0,
			LearnIterations:     0,
		},
	}

	attackInstance := CreateAttackInstance(&simulationInstance, attackSettings)
	if attackInstance == nil {
		panic("The requested attack type was not found. - " + attackSettings.AttackType)
	}

	startInstance := simulationInstance.DeepCopy() //We store a copy of the initial state to store it later in the database
	trackedState.UpdateSnapshot(startInstance)
	trackedState.StartTime = time.Now()
	syncReached := tpm_core.CompareWeights(settings.H, settings.K, settings.N,
		simulationInstance.StateA.Weights, simulationInstance.StateB.Weights)

	checkResult := 0
	skipCounter := trackConfig.SkipIterations
	for !syncReached {
		if simulationInstance.StimulateIterations > simConfig.IterationLimit {
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
				skipCounter = trackConfig.SkipIterations
			}
			skipCounter--
		}

	}

	//Get the best attackers, they will be saved in the db
	SetAttackerWeightOverlaps(attackSettings, *attackInstance)
	topAttackers := GetTopAttackersByWeight(*attackInstance)
	topScores := make([]float64, simConfig.StoreTopAttackerLimit)
	bestAttackerOverlap := -1.0
	bestAttackerScore := -1.0
	for i := range simConfig.StoreTopAttackerLimit {
		topScores[i] = topAttackers[i].weightOverlap
		if topAttackers[i].attackerScore > bestAttackerScore {
			bestAttackerScore = topAttackers[i].attackerScore
			bestAttackerOverlap = topAttackers[i].weightOverlap
		}
	}

	status := "LIMIT_REACHED"

	switch checkResult {
	case 0:
		status = "LIMIT_REACHED"
	case 1:
		if bestAttackerScore > ATTACKER_SCORE_THRESHOLD {
			status = "ATTACK_PARTIAL_SUCCESS"
		} else {
			status = "ON_SYNC"
		}
	case -1:
		status = "ATTACK_SUCCESS"
	case -2:
		status = "ATTACK_PARTIAL_SUCCESS"
	default:
		status = fmt.Sprintf("UNKOWN STATUS CODE: %d", checkResult)
		log.Fatalf("[ERROR]:  UNKOWN STATUS CODE: %d", checkResult)
	}

	result := AttackResult{
		Settings:            attackSettings,
		FinalState:          simulationInstance,
		SessionStatus:       status,
		StartTime:           trackedState.StartTime,
		EndTime:             time.Now(),
		AttackerOverlaps:    topScores,
		BestAttackerScore:   bestAttackerScore,
		BestAttackerOverlap: bestAttackerOverlap,
	}

	return result
}

// Helper function - used for tests only
func StimulateAllAttackers(settings AttackSettings, instance *AttackInstance, input_stimulus [][]int) {
	for _, v := range instance.attackerStates {
		v.Stimulate(settings.MTPMSettings, input_stimulus)
	}
}

func SetAttackerWeightOverlaps(settings AttackSettings, instance AttackInstance) {
	for _, attacker := range instance.attackerStates {
		overlap, score := tpm_core.CosineSimWeights(settings.H, settings.K, settings.N, instance.StateA.Weights, attacker.Weights)
		dataSize := engine.GetDataSize(settings.MTPMSettings)
		attacker.weightOverlap = overlap
		attacker.attackerScore = float64(score) / float64(dataSize)
	}
}

func GetTopAttackersByWeight(instance AttackInstance) []*AttackerState {
	// Make a copy of the slice of pointers
	attackers := make([]*AttackerState, len(instance.attackerStates))
	copy(attackers, instance.attackerStates)

	// Sort in-place by weightScore (highest first)
	sort.Slice(attackers, func(i, j int) bool {
		return attackers[i].weightOverlap > attackers[j].weightOverlap
	})

	return attackers
}

func CreateAttackInstance(simulationInstance *engine.SimulationInstance, attackSettings AttackSettings) *AttackInstance {
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
