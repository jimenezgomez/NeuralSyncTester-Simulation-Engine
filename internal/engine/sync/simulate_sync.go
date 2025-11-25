package sync

import (
	"time"

	config_manager "github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/config_manager/load"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/internal/engine"
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/pkg/tpm/tpm_core"
)

func SimulateSimpleSync(settings engine.MTPMSettings) engine.SimulationInstance {
	max_iterations := 100_000

	simulationInstance := engine.SimulationInstance{
		SimulationState: engine.SimulationState{
			StateA: engine.NewMTPMState(settings),
			StateB: engine.NewMTPMState(settings)},
		SimulationProgress: engine.SimulationProgress{
			StimulateIterations: 0,
			LearnIterations:     0,
		},
	}

	syncReached := tpm_core.CompareWeights(settings.H, settings.K, settings.N,
		simulationInstance.StateA.Weights, simulationInstance.StateB.Weights)

	for !syncReached {
		if simulationInstance.StimulateIterations > max_iterations {
			break
		}

		input_stimulus := tpm_core.CreateRandomStimulusArray(settings.K[0], settings.N[0], settings.M)
		simulationInstance.StateA.Stimulate(settings, input_stimulus)
		simulationInstance.StateB.Stimulate(settings, input_stimulus)
		simulationInstance.StimulateIterations += 1

		if simulationInstance.StateA.NetworkOutput == simulationInstance.StateB.NetworkOutput {
			simulationInstance.StateA.Learn(settings, simulationInstance.StateB.NetworkOutput)
			simulationInstance.StateB.Learn(settings, simulationInstance.StateA.NetworkOutput)

			simulationInstance.LearnIterations += 1
		}

		syncReached = engine.CompareWeights(settings, simulationInstance.StateA, simulationInstance.StateB)
	}
	return simulationInstance
}

func SimulateTrackedSync(trackedState *engine.TrackedMTPMSession, simConfig config_manager.SimulationConfig, trackConfig config_manager.TrackingConfig) engine.SimulationResult {
	settings := trackedState.Settings

	simulationInstance := engine.SimulationInstance{
		SimulationState: engine.SimulationState{
			StateA: engine.NewMTPMState(settings),
			StateB: engine.NewMTPMState(settings)},
		SimulationProgress: engine.SimulationProgress{
			StimulateIterations: 0,
			LearnIterations:     0,
		},
	}
	startInstance := simulationInstance.DeepCopy()
	trackedState.UpdateSnapshot(startInstance)
	trackedState.StartTime = time.Now()
	syncReached := tpm_core.CompareWeights(settings.H, settings.K, settings.N,
		simulationInstance.StateA.Weights, simulationInstance.StateB.Weights)

	skipCounter := trackConfig.SkipIterations
	for !syncReached {
		if simulationInstance.StimulateIterations > simConfig.IterationLimit {
			break
		}

		input_stimulus := tpm_core.CreateRandomStimulusArray(settings.K[0], settings.N[0], settings.M)
		simulationInstance.StateA.Stimulate(settings, input_stimulus)
		simulationInstance.StateB.Stimulate(settings, input_stimulus)
		simulationInstance.StimulateIterations += 1

		if simulationInstance.StateA.NetworkOutput == simulationInstance.StateB.NetworkOutput {
			simulationInstance.StateA.Learn(settings, simulationInstance.StateB.NetworkOutput)
			simulationInstance.StateB.Learn(settings, simulationInstance.StateA.NetworkOutput)

			simulationInstance.LearnIterations += 1
		}
		syncReached = tpm_core.CompareWeights(settings.H, settings.K, settings.N,
			simulationInstance.StateA.Weights, simulationInstance.StateB.Weights)

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

	status := "LIMIT_REACHED"
	if syncReached {
		status = "ON_SYNC"
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
