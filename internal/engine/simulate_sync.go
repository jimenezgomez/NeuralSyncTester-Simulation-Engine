package engine

import (
	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/pkg/tpm/tpm_core"
)

func SimulateSimpleSync(settings MTPMSettings) SimulationInstance {
	max_iterations := 100_000

	simulationInstance := SimulationInstance{
		StateA:              NewMTPMState(settings),
		StateB:              NewMTPMState(settings),
		StimulateIterations: 0,
		LearnIterations:     0,
	}

	syncReached := tpm_core.CompareWeights(settings.H, settings.K, settings.N,
		simulationInstance.StateA.Weights, simulationInstance.StateB.Weights)

	for !syncReached {
		if simulationInstance.StimulateIterations > max_iterations {
			break
		}

		input_stimulus := tpm_core.CreateRandomStimulusArray(settings.K[0], settings.N[0], settings.M)
		stimulate(settings, &simulationInstance.StateA, input_stimulus)
		stimulate(settings, &simulationInstance.StateB, input_stimulus)
		simulationInstance.StimulateIterations += 1

		if simulationInstance.StateA.NetworkOutput == simulationInstance.StateB.NetworkOutput {
			learn(settings, &simulationInstance.StateA, simulationInstance.StateB.NetworkOutput)
			learn(settings, &simulationInstance.StateB, simulationInstance.StateA.NetworkOutput)
			simulationInstance.LearnIterations += 1
		}

		syncReached = tpm_core.CompareWeights(settings.H, settings.K, settings.N,
			simulationInstance.StateA.Weights, simulationInstance.StateB.Weights)
	}
	return simulationInstance
}

func SimulateTrackedSync(trackedState *TrackedMTPMState) SimulationInstance {
	max_iterations := 1_000_000
	skipIterations := 150
	settings := trackedState.settings

	simulationInstance := SimulationInstance{
		StateA:              NewMTPMState(settings),
		StateB:              NewMTPMState(settings),
		StimulateIterations: 0,
		LearnIterations:     0,
	}
	snapshot := simulationInstance.DeepCopy()
	trackedState.UpdateSnapshot(snapshot)

	syncReached := tpm_core.CompareWeights(settings.H, settings.K, settings.N,
		simulationInstance.StateA.Weights, simulationInstance.StateB.Weights)

	skipCounter := skipIterations
	for !syncReached {
		if simulationInstance.StimulateIterations > max_iterations {
			break
		}

		input_stimulus := tpm_core.CreateRandomStimulusArray(settings.K[0], settings.N[0], settings.M)
		stimulate(settings, &simulationInstance.StateA, input_stimulus)
		stimulate(settings, &simulationInstance.StateB, input_stimulus)
		simulationInstance.StimulateIterations += 1

		if simulationInstance.StateA.NetworkOutput == simulationInstance.StateB.NetworkOutput {
			learn(settings, &simulationInstance.StateA, simulationInstance.StateB.NetworkOutput)
			learn(settings, &simulationInstance.StateB, simulationInstance.StateA.NetworkOutput)
			simulationInstance.LearnIterations += 1
		}

		syncReached = tpm_core.CompareWeights(settings.H, settings.K, settings.N,
			simulationInstance.StateA.Weights, simulationInstance.StateB.Weights)

		//TRACKING
		if trackedState.subCount.Load() > 0 {
			// Only bother tracking/publishing if someone is listening
			if skipCounter == 0 {
				snapshot := simulationInstance.DeepCopy()
				trackedState.UpdateSnapshot(snapshot)
				skipCounter = skipIterations
			}
			skipCounter--
		}

	}
	return simulationInstance
}
