package tpm_core

import "math"

func StimulateLayer(stimu [][]int, weights [][]int, k int, n int) []int {

	layerOutputs := make([]int, k)
	for i := 0; i < k; i++ {
		localField := NeuronLocalField(n, weights[i], stimu[i])
		localOutput := OutputSigma(localField)
		layerOutputs[i] = localOutput
	}

	return layerOutputs
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

func CosineSimWeights(h int, k []int, n []int, weights_a, weights_b [][][]int) float64 {
	dot := 0
	normA := 0
	normB := 0

	for layer := 0; layer < h; layer++ {
		for i := 0; i < k[layer]; i++ {
			for j := 0; j < n[layer]; j++ {
				va := weights_a[layer][i][j]
				vb := weights_b[layer][i][j]

				dot += va * vb
				normA += va * va
				normB += vb * vb
			}
		}
	}

	if normA == 0 || normB == 0 {
		return 0 // avoid division by zero
	}

	return float64(dot) / (math.Sqrt(float64(normA)) * math.Sqrt(float64(normB)))
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
