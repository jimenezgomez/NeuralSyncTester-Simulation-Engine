package config_manager

type Combination struct {
	K []int
	N []int
}

// buildK reconstructs the full K array from Kh and the N array.
// K[H-1] = Kh, and K[i] = K[i+1] * N[i+1] for i < H-1.
func buildK(N []int, Kh int) []int {
	H := len(N)
	K := make([]int, H)
	K[H-1] = Kh
	for i := H - 2; i >= 0; i-- {
		K[i] = K[i+1] * N[i+1]
	}
	return K
}

// buildNFullOverlap reconstructs the full N array from N0 and the K array.
// N[i] = K[i-1] for i < H-1.
func buildNFullOverlap(K []int, N0 int) []int {
	H := len(K)
	N := make([]int, H)
	N[0] = N0
	for i := 1; i < H; i++ {
		N[i] = K[i-1]
	}
	return N
}

// buildNPartialOverlap reconstructs the full N array from N0 and the K array.
// N[i] = K[i-1] - K[i] + 1 for i < H-1.
func buildNPartialOverlap(K []int, N0 int) []int {
	H := len(K)
	N := make([]int, H)
	N[0] = N0
	for i := 1; i < H; i++ {
		N[i] = K[i-1] - K[i] + 1
	}
	return N
}
