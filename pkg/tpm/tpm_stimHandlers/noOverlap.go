package tpm_stimHandlers

type NoOverlapTPM struct{}

func (tpm NoOverlapTPM) CreateStimulationStructure(n []int, k_last int) []int {
	h := len(n)
	k := make([]int, h)
	k[h-1] = k_last
	for i := 1; i < h; i++ {
		k[h-1-i] = n[h-i] * k[h-i]
	}
	return k
}

func (tpm NoOverlapTPM) CreateStimulusFromLayerOutput(dst [][]int, outputs []int, k_h int, n_h int) [][]int {
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
			dst[i][j] = outputs[n_h*i+j]
		}
	}

	return dst
}

func IntPow(base, exp int) int {
	result := 1
	for {
		if exp&1 == 1 {
			result *= base
		}
		exp >>= 1
		if exp == 0 {
			break
		}
		base *= base
	}

	return result
}
