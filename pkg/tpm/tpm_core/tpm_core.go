package tpm_core

import "math"

// StimulateLayer writes this layer's outputs into dst, reusing its backing
// array when it already has capacity k (it always does after the first call
// for a given layer, since k is fixed for the lifetime of a simulation).
func StimulateLayer(dst []int, stimu [][]int, weights [][]int, k int, n int) []int {
	if cap(dst) < k {
		dst = make([]int, k)
	} else {
		dst = dst[:k]
	}

	for i := 0; i < k; i++ {
		localField := NeuronLocalField(n, weights[i], stimu[i])
		localOutput := OutputSigma(localField)
		dst[i] = localOutput
	}

	return dst
}

func CompareWeights(h int, k []int, n []int, weights_a, weights_b [][][]int) bool {
	for layer := 0; layer < h; layer++ {
		for i := 0; i < k[layer]; i++ {
			for j := 0; j < n[layer]; j++ {
				if weights_a[layer][i][j] != weights_b[layer][i][j] {
					return false
				}
			}
		}
	}

	return true
}

func DotProdWeights(h int, k []int, n []int, weights_a, weights_b [][][]int) int {
	sum := 0
	for layer := 0; layer < h; layer++ {
		for i := 0; i < k[layer]; i++ {
			for j := 0; j < n[layer]; j++ {
				sum += weights_a[layer][i][j] * weights_b[layer][i][j]
			}
		}
	}

	return sum
}

// Calculates the overlap using the cosine similarity and gets the attacker scores as the amounr of correct values for each independent weight
func CosineSimWeights(h int, k []int, n []int, weights_a, weights_b [][][]int) (float64, int) {
	dot := 0
	normA := 0
	normB := 0

	score := 0

	for layer := 0; layer < h; layer++ {
		for i := 0; i < k[layer]; i++ {
			for j := 0; j < n[layer]; j++ {
				va := weights_a[layer][i][j]
				vb := weights_b[layer][i][j]

				dot += va * vb
				normA += va * va
				normB += vb * vb

				// If these weights are the same value at the same position, add score
				if va == vb {
					score++
				}

			}
		}
	}

	if normA == 0 || normB == 0 {
		return 0, score // avoid division by zero
	}

	return float64(dot) / (math.Sqrt(float64(normA)) * math.Sqrt(float64(normB))), score
}

func CreateRandomStimulusArray(k int, n int, m int) [][]int {
	stim := make([][]int, k)
	for i := 0; i < k; i++ {
		stim[i] = make([]int, n)
		for j := 0; j < n; j++ {
			stim[i][j] = (CryptoRandIntn(2)*2 - 1) * (CryptoRandIntn(m) + 1)
		}
	}
	return stim
}

func CreateRandomLayerWeightsArray(k int, n int, l int) [][]int {
	w := make([][]int, k)
	for i := 0; i < k; i++ {
		w[i] = make([]int, n)
		for j := 0; j < n; j++ {
			w[i][j] = (CryptoRandIntn(2)*2 - 1) * (CryptoRandIntn(l + 1)) // l + 1 because the function goes from [0,l[
		}
	}
	return w
}

func GetNetworkDataSize(H int, K []int, N []int) int {
	totalDataSize := 0
	for layer := 0; layer < H; layer++ {
		totalDataSize += K[layer] * N[layer]
	}
	return totalDataSize
}
