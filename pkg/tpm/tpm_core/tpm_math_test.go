package tpm_core_test

import (
	"math"
	"testing"

	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/pkg/tpm/tpm_core"
)

func TestNeuronLocalFieldRaw(t *testing.T) {
	tests := []struct {
		name string
		n    int
		w    []int
		stim []int
		want float64
	}{
		{"zero vectors", 3, []int{0, 0, 0}, []int{0, 0, 0}, 0},
		{"simple positive", 3, []int{1, 2, 3}, []int{1, 1, 1}, 6},
		{"mixed sign", 3, []int{1, -2, 3}, []int{1, 1, -1}, -4},
		{"n truncates slice", 2, []int{1, 2, 3}, []int{1, 1, 1}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tpm_core.NeuronLocalFieldRaw(tt.n, tt.w, tt.stim)
			if got != tt.want {
				t.Errorf("NeuronLocalFieldRaw() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNeuronLocalField(t *testing.T) {
	tests := []struct {
		name string
		n    int
		w    []int
		stim []int
	}{
		{"n=1", 1, []int{3}, []int{2}},
		{"n=4", 4, []int{1, 2, 3, 4}, []int{1, 1, 1, 1}},
		{"n=16", 16, []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := tpm_core.NeuronLocalFieldRaw(tt.n, tt.w, tt.stim)
			want := raw * (1 / math.Sqrt(float64(tt.n)))
			got := tpm_core.NeuronLocalField(tt.n, tt.w, tt.stim)
			// FastInverseSqrt is a one-Newton-iteration approximation, tolerate ~0.2% relative error.
			epsilon := math.Abs(want)*0.002 + 1e-9
			if !floatsAlmostEqual(got, want, epsilon) {
				t.Errorf("NeuronLocalField() = %v, want ~%v (epsilon %v)", got, want, epsilon)
			}
		})
	}

	t.Run("raw=0 is exactly 0 regardless of n", func(t *testing.T) {
		got := tpm_core.NeuronLocalField(2, []int{0, 0}, []int{5, -5})
		if got != 0 {
			t.Errorf("NeuronLocalField() = %v, want 0", got)
		}
	})
}

func TestOutputSigma(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		want int
	}{
		{"positive", 5, 1},
		{"negative", -5, -1},
		{"zero maps to -1", 0, -1},
		{"tiny negative", -0.0001, -1},
		{"tiny positive", math.SmallestNonzeroFloat64, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tpm_core.OutputSigma(tt.x); got != tt.want {
				t.Errorf("OutputSigma(%v) = %v, want %v", tt.x, got, tt.want)
			}
		})
	}
}

func TestThau(t *testing.T) {
	tests := []struct {
		name    string
		outputs []int
		k       int
		want    int
	}{
		{"all agree positive", []int{1, 1, 1}, 3, 1},
		{"one disagrees", []int{1, -1, 1}, 3, -1},
		{"two negatives multiply positive", []int{-1, -1}, 2, 1},
		{"three negatives multiply negative", []int{-1, -1, -1}, 3, -1},
		{"k=0 empty product identity", []int{-1, -1, -1}, 0, 1},
		{"k shorter than outputs", []int{-1, -1, 1, 1}, 2, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tpm_core.Thau(tt.outputs, tt.k); got != tt.want {
				t.Errorf("Thau() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHeavisideStep(t *testing.T) {
	tests := []struct {
		x    int
		want int
	}{
		{5, 1},
		{0, 0},
		{-5, 0},
	}
	for _, tt := range tests {
		if got := tpm_core.HeavisideStep(tt.x); got != tt.want {
			t.Errorf("HeavisideStep(%v) = %v, want %v", tt.x, got, tt.want)
		}
	}
}

func TestGFunction(t *testing.T) {
	tests := []struct {
		name string
		w, l int
		want int
	}{
		{"no clip zero", 0, 5, 0},
		{"no clip positive", 3, 5, 3},
		{"exact boundary positive", 5, 5, 5},
		{"over boundary positive", 6, 5, 5},
		{"no clip negative", -3, 5, -3},
		{"over boundary negative", -6, 5, -5},
		{"l=0 forces zero positive", 3, 0, 0},
		{"l=0 forces zero negative", -3, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tpm_core.GFunction(tt.w, tt.l); got != tt.want {
				t.Errorf("GFunction(%v, %v) = %v, want %v", tt.w, tt.l, got, tt.want)
			}
		})
	}
}

func TestFastInverseSqrt(t *testing.T) {
	tests := []float64{1, 4, 16, 1000, 0.0001}
	for _, x := range tests {
		t.Run("", func(t *testing.T) {
			want := 1 / math.Sqrt(x)
			got := tpm_core.FastInverseSqrt(x)
			epsilon := want * 0.002
			if !floatsAlmostEqual(got, want, epsilon) {
				t.Errorf("FastInverseSqrt(%v) = %v, want ~%v (epsilon %v)", x, got, want, epsilon)
			}
		})
	}

	// Documented edge cases: the bit-trick doesn't special-case these inputs, so the
	// result isn't meaningful, but the function must not panic.
	t.Run("ZeroInput_DocumentedEdgeCase", func(t *testing.T) {
		got := tpm_core.FastInverseSqrt(0)
		if math.IsNaN(got) {
			t.Errorf("FastInverseSqrt(0) = NaN, expected a defined (if meaningless) float")
		}
	})

	t.Run("NegativeInput_UndefinedBehavior", func(t *testing.T) {
		// The sign bit corrupts the magic-number subtraction; just assert no panic occurs.
		_ = tpm_core.FastInverseSqrt(-4)
	})
}

func TestCryptoRandIntnErr(t *testing.T) {
	t.Run("n=0 returns error", func(t *testing.T) {
		got, err := tpm_core.CryptoRandIntn_err(0)
		if err == nil || got != 0 {
			t.Errorf("CryptoRandIntn_err(0) = (%v, %v), want (0, non-nil error)", got, err)
		}
	})

	t.Run("negative n returns error", func(t *testing.T) {
		got, err := tpm_core.CryptoRandIntn_err(-5)
		if err == nil || got != 0 {
			t.Errorf("CryptoRandIntn_err(-5) = (%v, %v), want (0, non-nil error)", got, err)
		}
	})

	t.Run("n=1 always returns 0", func(t *testing.T) {
		for i := 0; i < 200; i++ {
			got, err := tpm_core.CryptoRandIntn_err(1)
			if err != nil || got != 0 {
				t.Fatalf("CryptoRandIntn_err(1) = (%v, %v), want (0, nil)", got, err)
			}
		}
	})

	t.Run("n=10 stays in range", func(t *testing.T) {
		for i := 0; i < 2000; i++ {
			got, err := tpm_core.CryptoRandIntn_err(10)
			if err != nil {
				t.Fatalf("CryptoRandIntn_err(10) returned error: %v", err)
			}
			if got < 0 || got >= 10 {
				t.Fatalf("CryptoRandIntn_err(10) = %v, want in [0,10)", got)
			}
		}
	})
}

func TestCryptoRandIntn(t *testing.T) {
	t.Run("n=0 returns 0", func(t *testing.T) {
		if got := tpm_core.CryptoRandIntn(0); got != 0 {
			t.Errorf("CryptoRandIntn(0) = %v, want 0", got)
		}
	})

	t.Run("negative n returns 0", func(t *testing.T) {
		if got := tpm_core.CryptoRandIntn(-3); got != 0 {
			t.Errorf("CryptoRandIntn(-3) = %v, want 0", got)
		}
	})

	t.Run("n=1 always returns 0", func(t *testing.T) {
		for i := 0; i < 200; i++ {
			if got := tpm_core.CryptoRandIntn(1); got != 0 {
				t.Fatalf("CryptoRandIntn(1) = %v, want 0", got)
			}
		}
	})

	t.Run("n=5 stays in range", func(t *testing.T) {
		for i := 0; i < 2000; i++ {
			got := tpm_core.CryptoRandIntn(5)
			if got < 0 || got >= 5 {
				t.Fatalf("CryptoRandIntn(5) = %v, want in [0,5)", got)
			}
		}
	})
}
