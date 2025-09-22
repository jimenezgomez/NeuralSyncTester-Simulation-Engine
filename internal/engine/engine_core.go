package engine

import "github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/pkg/tpm/tpm_core"

func stimulate(settings MTPMSettings, sessionState *MTPMState, firstLayerInput [][]int) {
	inputs := sessionState.InputBuffer
	inputs[0] = firstLayerInput
	outputs := sessionState.OutputBuffer

	//Stimulate all layers, to avoid overflowing the inputs array we do the last layer separately -> (There is no next layer, no more inputs)
	for layer := 0; layer < settings.H-1; layer++ {
		outputs[layer] = tpm_core.StimulateLayer(inputs[layer], sessionState.Weights[layer], settings.K[layer], settings.N[layer])
		inputs[layer+1] = settings.stimulationHandlers.CreateStimulusFromLayerOutput(outputs[layer], settings.K[layer+1], settings.N[layer+1])
	}
	outputs[settings.H-1] = tpm_core.StimulateLayer(inputs[settings.H-1], sessionState.Weights[settings.H-1], settings.K[settings.H-1], settings.N[settings.H-1])
	sessionState.NetworkOutput = tpm_core.Thau(outputs[settings.H-1], settings.K[settings.H-1])
}

func learn(settings MTPMSettings, mtpmState *MTPMState, remoteOutput int) {
	for layer := 0; layer < settings.H; layer++ {
		settings.learnRuleHandler.TPMLearnLayer(
			settings.K[layer], settings.N[layer], settings.L,
			mtpmState.Weights[layer],
			mtpmState.InputBuffer[layer], mtpmState.OutputBuffer[layer],
			mtpmState.NetworkOutput, remoteOutput)
	}
}
func GetDataSize(settings MTPMSettings) int {
	return tpm_core.GetNetworkDataSize(settings.H, settings.K, settings.N)
}
