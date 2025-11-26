package tpm_core

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
)

func NeuronLocalFieldRaw(n int, w_k []int, stim_k []int) float64 {
	dot_prod := 0
	for i := 0; i < n; i++ {
		dot_prod += w_k[i] * stim_k[i]
	}
	return float64(dot_prod)
}

func NeuronLocalField(n int, w_k []int, stim_k []int) float64 {
	return NeuronLocalFieldRaw(n, w_k, stim_k) * FastInverseSqrt(float64(n))
}

func OutputSigma(x float64) int {
	if x > 0 {
		return 1
	}
	return -1
}

func Thau(outputs []int, k int) int {
	mul := 1
	for i := 0; i < k; i++ {
		mul *= outputs[i]
	}

	return mul
}

func HeavisideStep(x int) int {
	if x > 0 {
		return 1
	}
	return 0
}

func GFunction(w int, l int) int {
	sign := 1
	if w < 0 {
		sign = -1
	}
	if w*sign > l {
		return l * sign
	}
	return w
}

func FastInverseSqrt(x float64) float64 {
	i := math.Float64bits(x)
	i = 0x5fe6eb50c7b537a9 - (i >> 1)
	y := math.Float64frombits(i)

	y = y * (1.5 - 0.5*x*y*y)
	return y
}

// CryptoRandIntn returns a cryptographically secure pseudo-random number in [0, n).
// It panics if n <= 0.
// It is safe for concurrent use.
func CryptoRandIntn_err(n int) (int, error) {
	if n <= 0 {
		return 0, fmt.Errorf("argument to CryptoRandIntn must be positive, got %d", n)
	}

	// big.Int is necessary because crypto/rand.Int operates on arbitrary-precision integers
	// to handle potentially very large ranges without overflow issues.
	// We convert the input 'n' to a big.Int.
	max := big.NewInt(int64(n))

	// crypto/rand.Int generates a cryptographically secure random number in [0, max).
	// It efficiently avoids modulo bias.
	result, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, fmt.Errorf("failed to generate crypto random number: %w", err)
	}

	// Convert the big.Int result back to an int.
	return int(result.Int64()), nil
}

// CryptoRandIntn returns a cryptographically secure pseudo-random number in [0, n)
// This function is the same, but for direct replacement into current code
func CryptoRandIntn(n int) int {
	if n <= 0 {
		return 0
	}

	// big.Int is necessary because crypto/rand.Int operates on arbitrary-precision integers
	// to handle potentially very large ranges without overflow issues.
	// We convert the input 'n' to a big.Int.
	max := big.NewInt(int64(n))

	// crypto/rand.Int generates a cryptographically secure random number in [0, max).
	// It efficiently avoids modulo bias.
	result, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0
	}

	// Convert the big.Int result back to an int.
	return int(result.Int64())
}
