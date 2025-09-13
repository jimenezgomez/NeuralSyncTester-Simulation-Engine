package tpm_core

import (
	"math"
	"math/rand"
)

func StimulateLayer(stimu [][]int, weights [][]int, k int, n int) []int {

	layerOutputs := make([]int, k)
	for i := 0; i < k; i++ {
		localField := NeuronLocalField(n, weights[i], stimu[i])
		localOutput := OutputSigma(localField)
		layerOutputs[i] = localOutput
	}

	return layerOutputs
}

func CompareWeights(h int, k []int, n []int, weights_a [][][]int, weights_b [][][]int) bool {
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

func CreateRandomStimulusArray(k int, n int, m int, localRand *rand.Rand) [][]int {
	stim := make([][]int, k)
	for i := 0; i < k; i++ {
		stim[i] = make([]int, n)
		for j := 0; j < n; j++ {
			stim[i][j] = (localRand.Intn(2)*2 - 1) * (localRand.Intn(m) + 1)
		}
	}
	return stim
}

func CreateRandomLayerWeightsArray(k int, n int, l int, localRand *rand.Rand) [][]int {
	w := make([][]int, k)
	for i := 0; i < k; i++ {
		w[i] = make([]int, n)
		for j := 0; j < n; j++ {
			w[i][j] = (localRand.Intn(2)*2 - 1) * (localRand.Intn(l + 1)) // l + 1 because the function goes from [0,l[
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
