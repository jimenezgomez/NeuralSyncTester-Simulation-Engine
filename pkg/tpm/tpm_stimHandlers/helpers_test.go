package tpm_stimHandlers_test

import "github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/pkg/tpm/tpm_stimHandlers"

// seqInts builds sentinel values [start, start+1, ..., start+n-1] so index-mapping
// bugs in the handlers can't hide behind repeated ±1 output values.
func seqInts(start, n int) []int {
	out := make([]int, n)
	for i := 0; i < n; i++ {
		out[i] = start + i
	}
	return out
}

var (
	_ tpm_stimHandlers.TPMStimulationHandlers = tpm_stimHandlers.FullOverlapTPM{}
	_ tpm_stimHandlers.TPMStimulationHandlers = tpm_stimHandlers.NoOverlapTPM{}
	_ tpm_stimHandlers.TPMStimulationHandlers = tpm_stimHandlers.PartialOverlapTPM{}
)
