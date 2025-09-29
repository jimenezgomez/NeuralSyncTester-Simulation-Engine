package engine

import "github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/pkg/tpm/tpm_core"

func (mtpmState *MTPMState) Stimulate(settings MTPMSettings, firstLayerInput [][]int) {
	inputs := mtpmState.InputBuffer
	inputs[0] = firstLayerInput
	outputs := mtpmState.OutputBuffer

	//Stimulate all layers, to avoid overflowing the inputs array we do the last layer separately -> (There is no next layer, no more inputs)
	for layer := 0; layer < settings.H-1; layer++ {
		//calculate outputs for this layer
		outputs[layer] = tpm_core.StimulateLayer(
			inputs[layer], mtpmState.Weights[layer],
			settings.K[layer], settings.N[layer])
		//set inputs for next layer
		inputs[layer+1] = settings.stimulationHandlers.CreateStimulusFromLayerOutput(
			outputs[layer],
			settings.K[layer+1], settings.N[layer+1])
	}
	outputs[settings.H-1] = tpm_core.StimulateLayer(
		inputs[settings.H-1], mtpmState.Weights[settings.H-1],
		settings.K[settings.H-1], settings.N[settings.H-1])
	mtpmState.NetworkOutput = tpm_core.Thau(outputs[settings.H-1], settings.K[settings.H-1])
}

func (mtpmState *MTPMState) Learn(settings MTPMSettings, remoteOutput int) {
	for layer := 0; layer < settings.H; layer++ {
		settings.learnRuleHandler.TPMLearnLayer(
			settings.K[layer], settings.N[layer], settings.L,
			mtpmState.Weights[layer],
			mtpmState.InputBuffer[layer], mtpmState.OutputBuffer[layer],
			mtpmState.NetworkOutput, remoteOutput)
	}
}

func (mtpmState *MTPMState) LearnWithOutputs(settings MTPMSettings, output_A, output_B int) {
	for layer := 0; layer < settings.H; layer++ {
		settings.learnRuleHandler.TPMLearnLayer(
			settings.K[layer], settings.N[layer], settings.L,
			mtpmState.Weights[layer],
			mtpmState.InputBuffer[layer], mtpmState.OutputBuffer[layer],
			output_A, output_B)
	}
}

func (mtpmState *MTPMState) LearnWithFullReference(settings MTPMSettings, output_A, output_B int, referenceOutputBuffer *[][]int) {
	for layer := 0; layer < settings.H; layer++ {
		settings.learnRuleHandler.TPMLearnLayer(
			settings.K[layer], settings.N[layer], settings.L,
			mtpmState.Weights[layer],
			mtpmState.InputBuffer[layer], (*referenceOutputBuffer)[layer],
			output_A, output_B)
	}
}

func GetDataSize(settings MTPMSettings) int {
	return tpm_core.GetNetworkDataSize(settings.H, settings.K, settings.N)
}

func CompareWeights(settings MTPMSettings, stateA, stateB MTPMState) bool {
	return tpm_core.CompareWeights(settings.H, settings.K, settings.N,
		stateA.Weights, stateB.Weights)
}
