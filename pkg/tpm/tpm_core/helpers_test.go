package tpm_core_test

import "math"

func floatsAlmostEqual(got, want, epsilon float64) bool {
	return math.Abs(got-want) <= epsilon
}

// buildWeights3D builds an h-layer weight array with layer i having k[i] neurons
// each with n[i] weights, values produced by gen(layer, i, j).
func buildWeights3D(h int, k, n []int, gen func(layer, i, j int) int) [][][]int {
	w := make([][][]int, h)
	for layer := 0; layer < h; layer++ {
		w[layer] = make([][]int, k[layer])
		for i := 0; i < k[layer]; i++ {
			w[layer][i] = make([]int, n[layer])
			for j := 0; j < n[layer]; j++ {
				w[layer][i][j] = gen(layer, i, j)
			}
		}
	}
	return w
}

func allInRange(vals []int, lo, hi int) bool {
	for _, v := range vals {
		if v < lo || v > hi {
			return false
		}
	}
	return true
}
