package tpm_core

func NeuronLocalField(n int, w_k []int, stim_k []int) float64 {
	dot_prod := 0
	for i := 0; i < n; i++ {
		dot_prod += w_k[i] * stim_k[i]
	}
	return float64(dot_prod) * FastInverseSqrt(float64(n))
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