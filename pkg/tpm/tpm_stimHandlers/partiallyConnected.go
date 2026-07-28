package tpm_stimHandlers

type PartialOverlapTPM struct{}

func (tpm PartialOverlapTPM) CreateStimulationStructure(k []int, n_0 int) []int {
	prev := -1
	h := len(k)
	n := make([]int, h)

	n[0] = n_0
	prev = k[0]

	for layer := 1; layer < h; layer++ {

		if k[layer] >= prev {
			return nil
		}
		n[layer] = k[layer-1] - k[layer] + 1
		prev = k[layer]
	}
	return n
}

func (tpm PartialOverlapTPM) CreateStimulusFromLayerOutput(dst [][]int, outputs []int, k_h int, n_h int) [][]int {
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
		for j := 0; j < n_h; j++ {
			dst[i][j] = outputs[j+i]
		}
	}

	return dst
}
