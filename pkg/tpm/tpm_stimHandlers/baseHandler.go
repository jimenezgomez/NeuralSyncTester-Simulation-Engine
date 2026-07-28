package tpm_stimHandlers

type TPMStimulationHandlers interface {
	CreateStimulationStructure(k []int, n_0 int) []int
	CreateStimulusFromLayerOutput(dst [][]int, outputs []int, k_h int, n_h int) [][]int
}
