package tpm_stimHandlers

type FullOverlapTPM struct{}

func (tpm FullOverlapTPM) CreateStimulationStructure(k []int, n_0 int) []int {
	h := len(k)
	n := make([]int, h)

	n[0] = n_0

	for layer := 1; layer < h; layer++ {
		n[layer] = k[layer-1]
	}
	return n
}

func (tpm FullOverlapTPM) CreateStimulusFromLayerOutput(dst [][]int, outputs []int, k_h int, n_h int) [][]int {
	if cap(dst) < k_h {
		dst = make([][]int, k_h)
	} else {
		dst = dst[:k_h]
	}
	for i := 0; i < k_h; i++ {
		if cap(dst[i]) < n_h {
			dst[i] = make([]int, n_h)
		} else {
			dst[i] = dst[i][:n_h]
		}
		//When fully connected, the stim count is the same as the neuron count from the prev layer
		for j := 0; j < n_h; j++ {
			dst[i][j] = outputs[j] //So this maps outputs to inputs, 1 to 1
		}
	}

	return dst
}
