package attacks

import (
	"time"

	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/pkg/tpm/tpm_core"
)

func RunTrackedAttack_Simple(trackedState *engine.TrackedMTPMState) engine.SimulationResult {
	max_iterations := 100_000
	skipIterations := 150
	settings := trackedState.GetSettings()

	attackSettings := AttackSettings{
		MTPMSettings:  settings,
		attackerLimit: 10,
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

	attackInstance := NewSimpleAttack(attackSettings, simulationInstance)

	startInstance := simulationInstance.DeepCopy()
	trackedState.UpdateSnapshot(startInstance)
	trackedState.StartTime = time.Now()
	syncReached := tpm_core.CompareWeights(settings.H, settings.K, settings.N,
		simulationInstance.StateA.Weights, simulationInstance.StateB.Weights)

	checkResult := 0
	skipCounter := skipIterations
	for !syncReached {
		if simulationInstance.StimulateIterations > max_iterations {
			break
		}

		input_stimulus := tpm_core.CreateRandomStimulusArray(settings.K[0], settings.N[0], settings.M)
		simulationInstance.StateA.Stimulate(settings, input_stimulus)
		simulationInstance.StateB.Stimulate(settings, input_stimulus)
		// StimulateAllAttackers(attackSettings, attackInstance, input_stimulus)
		// Here we should stimulate the attackers, but we don't really need to unless we need to update their weights
		simulationInstance.StimulateIterations += 1

		if simulationInstance.StateA.NetworkOutput == simulationInstance.StateB.NetworkOutput {
			simulationInstance.StateA.Learn(settings, simulationInstance.StateB.NetworkOutput)
			simulationInstance.StateB.Learn(settings, simulationInstance.StateA.NetworkOutput)
			attackInstance.attackerExec(
				attackSettings, attackInstance,
				attackInstance.StateA.NetworkOutput, attackInstance.StateB.NetworkOutput,
				input_stimulus)
			simulationInstance.LearnIterations += 1
		}

		checkResult = attackInstance.attackerCheck(attackSettings, attackInstance)
		syncReached = checkResult != 0

		//TRACKING
		if trackedState.GetSubCount() > 0 {
			// Only bother tracking/publishing if someone is listening
			if skipCounter == 0 {
				snapshot := simulationInstance.DeepCopy()
				trackedState.UpdateSnapshot(snapshot)
				skipCounter = skipIterations
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

	result := engine.SimulationResult{
		Settings:      settings,
		InitialState:  startInstance.SimulationState,
		FinalState:    simulationInstance.SimulationState,
		Iterations:    simulationInstance.SimulationProgress,
		SessionStatus: status,
		StartTime:     trackedState.StartTime,
		EndTime:       time.Now(),
	}

	return result
}

// Helper function - used for tests only
func StimulateAllAttackers(settings AttackSettings, instance AttackInstance, input_stimulus [][]int) {
	for _, v := range instance.attackerStates {
		v.Stimulate(settings.MTPMSettings, input_stimulus)
	}
}

func RunTrackedAttack_Geom(trackedState *engine.TrackedMTPMState) engine.SimulationResult {
	max_iterations := 100_000
	skipIterations := 150
	settings := trackedState.GetSettings()

	attackSettings := AttackSettings{
		MTPMSettings:  settings,
		attackerLimit: 10,
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

	attackInstance := NewGeomAttack(attackSettings, simulationInstance)

	startInstance := simulationInstance.DeepCopy()
	trackedState.UpdateSnapshot(startInstance)
	trackedState.StartTime = time.Now()
	syncReached := tpm_core.CompareWeights(settings.H, settings.K, settings.N,
		simulationInstance.StateA.Weights, simulationInstance.StateB.Weights)

	checkResult := 0
	skipCounter := skipIterations
	for !syncReached {
		if simulationInstance.StimulateIterations > max_iterations {
			break
		}

		input_stimulus := tpm_core.CreateRandomStimulusArray(settings.K[0], settings.N[0], settings.M)
		simulationInstance.StateA.Stimulate(settings, input_stimulus)
		simulationInstance.StateB.Stimulate(settings, input_stimulus)
		// StimulateAllAttackers(attackSettings, attackInstance, input_stimulus)
		// Here we should stimulate the attackers, but we don't really need to unless we need to update their weights
		simulationInstance.StimulateIterations += 1

		if simulationInstance.StateA.NetworkOutput == simulationInstance.StateB.NetworkOutput {
			simulationInstance.StateA.Learn(settings, simulationInstance.StateB.NetworkOutput)
			simulationInstance.StateB.Learn(settings, simulationInstance.StateA.NetworkOutput)
			attackInstance.attackerExec(
				attackSettings, attackInstance,
				attackInstance.StateA.NetworkOutput, attackInstance.StateB.NetworkOutput,
				input_stimulus)
			simulationInstance.LearnIterations += 1
		}

		checkResult = attackInstance.attackerCheck(attackSettings, attackInstance)
		syncReached = checkResult != 0

		//TRACKING
		if trackedState.GetSubCount() > 0 {
			// Only bother tracking/publishing if someone is listening
			if skipCounter == 0 {
				snapshot := simulationInstance.DeepCopy()
				trackedState.UpdateSnapshot(snapshot)
				skipCounter = skipIterations
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

	result := engine.SimulationResult{
		Settings:      settings,
		InitialState:  startInstance.SimulationState,
		FinalState:    simulationInstance.SimulationState,
		Iterations:    simulationInstance.SimulationProgress,
		SessionStatus: status,
		StartTime:     trackedState.StartTime,
		EndTime:       time.Now(),
	}

	return result
}

func RunTrackedAttack_Maj(trackedState *engine.TrackedMTPMState) engine.SimulationResult {
	max_iterations := 100_000
	skipIterations := 150
	settings := trackedState.GetSettings()

	attackSettings := AttackSettings{
		MTPMSettings:  settings,
		attackerLimit: 10,
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

	attackInstance := NewMajorityAttack(attackSettings, simulationInstance)

	startInstance := simulationInstance.DeepCopy()
	trackedState.UpdateSnapshot(startInstance)
	trackedState.StartTime = time.Now()
	syncReached := tpm_core.CompareWeights(settings.H, settings.K, settings.N,
		simulationInstance.StateA.Weights, simulationInstance.StateB.Weights)

	checkResult := 0
	skipCounter := skipIterations
	for !syncReached {
		if simulationInstance.StimulateIterations > max_iterations {
			break
		}

		input_stimulus := tpm_core.CreateRandomStimulusArray(settings.K[0], settings.N[0], settings.M)
		simulationInstance.StateA.Stimulate(settings, input_stimulus)
		simulationInstance.StateB.Stimulate(settings, input_stimulus)
		// StimulateAllAttackers(attackSettings, attackInstance, input_stimulus)
		// Here we should stimulate the attackers, but we don't really need to unless we need to update their weights
		simulationInstance.StimulateIterations += 1

		if simulationInstance.StateA.NetworkOutput == simulationInstance.StateB.NetworkOutput {
			simulationInstance.StateA.Learn(settings, simulationInstance.StateB.NetworkOutput)
			simulationInstance.StateB.Learn(settings, simulationInstance.StateA.NetworkOutput)
			attackInstance.attackerExec(
				attackSettings, attackInstance,
				attackInstance.StateA.NetworkOutput, attackInstance.StateB.NetworkOutput,
				input_stimulus)
			simulationInstance.LearnIterations += 1
		}

		checkResult = attackInstance.attackerCheck(attackSettings, attackInstance)
		syncReached = checkResult != 0

		//TRACKING
		if trackedState.GetSubCount() > 0 {
			// Only bother tracking/publishing if someone is listening
			if skipCounter == 0 {
				snapshot := simulationInstance.DeepCopy()
				trackedState.UpdateSnapshot(snapshot)
				skipCounter = skipIterations
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

	result := engine.SimulationResult{
		Settings:      settings,
		InitialState:  startInstance.SimulationState,
		FinalState:    simulationInstance.SimulationState,
		Iterations:    simulationInstance.SimulationProgress,
		SessionStatus: status,
		StartTime:     trackedState.StartTime,
		EndTime:       time.Now(),
	}

	return result
}
