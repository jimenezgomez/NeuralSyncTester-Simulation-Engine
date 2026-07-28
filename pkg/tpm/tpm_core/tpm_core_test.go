package tpm_core_test

import (
	"testing"

	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/pkg/tpm/tpm_core"
)

func TestStimulateLayer(t *testing.T) {
	t.Run("hand-computed outputs including zero-field edge case", func(t *testing.T) {
		weights := [][]int{{1, 1}, {-1, -1}, {1, -1}}
		stimu := [][]int{{1, 1}, {1, 1}, {1, 1}}
		// local fields: neuron0=2 (>0->1), neuron1=-2 (<0->-1), neuron2=0 (OutputSigma(0)=-1)
		want := []int{1, -1, -1}

		got := tpm_core.StimulateLayer(nil, stimu, weights, 3, 2)
		if len(got) != len(want) {
			t.Fatalf("StimulateLayer() length = %v, want %v", len(got), len(want))
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("StimulateLayer()[%d] = %v, want %v", i, got[i], want[i])
			}
		}
	})

	t.Run("reuses dst backing array when capacity suffices", func(t *testing.T) {
		dst := make([]int, 5, 10)
		weights := [][]int{{1, 1}, {1, 1}, {1, 1}}
		stimu := [][]int{{1, 1}, {1, 1}, {1, 1}}

		got := tpm_core.StimulateLayer(dst, stimu, weights, 3, 2)
		if len(got) != 3 {
			t.Errorf("len(got) = %v, want 3", len(got))
		}
		if cap(got) != 10 {
			t.Errorf("cap(got) = %v, want 10 (backing array reused)", cap(got))
		}
	})

	t.Run("allocates fresh backing array when capacity insufficient", func(t *testing.T) {
		dst := make([]int, 0, 1)
		weights := [][]int{{1, 1}, {1, 1}, {1, 1}}
		stimu := [][]int{{1, 1}, {1, 1}, {1, 1}}

		got := tpm_core.StimulateLayer(dst, stimu, weights, 3, 2)
		if cap(got) != 3 {
			t.Errorf("cap(got) = %v, want 3 (fresh allocation)", cap(got))
		}
	})
}

func TestCompareWeights(t *testing.T) {
	k := []int{2, 1}
	n := []int{2, 3}
	gen := func(layer, i, j int) int { return layer*100 + i*10 + j }

	t.Run("identical arrays", func(t *testing.T) {
		a := buildWeights3D(2, k, n, gen)
		b := buildWeights3D(2, k, n, gen)
		if !tpm_core.CompareWeights(2, k, n, a, b) {
			t.Error("CompareWeights() = false, want true for identical arrays")
		}
	})

	t.Run("single mismatch inside bounds", func(t *testing.T) {
		a := buildWeights3D(2, k, n, gen)
		b := buildWeights3D(2, k, n, gen)
		b[1][0][2] = 9999
		if tpm_core.CompareWeights(2, k, n, a, b) {
			t.Error("CompareWeights() = true, want false when a value differs")
		}
	})

	t.Run("h=0 is vacuously true", func(t *testing.T) {
		a := buildWeights3D(2, k, n, gen)
		b := buildWeights3D(2, k, n, gen)
		b[0][0][0] = 9999
		if !tpm_core.CompareWeights(0, k, n, a, b) {
			t.Error("CompareWeights() = false, want true when h=0")
		}
	})

	t.Run("mismatch outside h/k/n bounds is ignored", func(t *testing.T) {
		a := buildWeights3D(2, k, n, gen)
		b := buildWeights3D(2, k, n, gen)
		// h only covers layer 0 here; a mismatch in layer 1 must not affect the result.
		b[1][0][2] = 9999
		if !tpm_core.CompareWeights(1, k, n, a, b) {
			t.Error("CompareWeights() = false, want true when mismatch is outside the h bound")
		}
	})
}

func TestDotProdWeights(t *testing.T) {
	k := []int{2}
	n := []int{2}

	t.Run("hand-computed", func(t *testing.T) {
		a := [][][]int{{{1, 2}, {3, 4}}}
		b := [][][]int{{{5, 6}, {7, 8}}}
		want := 1*5 + 2*6 + 3*7 + 4*8
		if got := tpm_core.DotProdWeights(1, k, n, a, b); got != want {
			t.Errorf("DotProdWeights() = %v, want %v", got, want)
		}
	})

	t.Run("either operand all zero", func(t *testing.T) {
		a := [][][]int{{{0, 0}, {0, 0}}}
		b := [][][]int{{{5, 6}, {7, 8}}}
		if got := tpm_core.DotProdWeights(1, k, n, a, b); got != 0 {
			t.Errorf("DotProdWeights() = %v, want 0", got)
		}
	})

	t.Run("h=0", func(t *testing.T) {
		a := [][][]int{{{1, 2}, {3, 4}}}
		b := [][][]int{{{5, 6}, {7, 8}}}
		if got := tpm_core.DotProdWeights(0, k, n, a, b); got != 0 {
			t.Errorf("DotProdWeights() = %v, want 0", got)
		}
	})
}

func TestCosineSimWeights(t *testing.T) {
	k := []int{1}
	n := []int{2}

	t.Run("identical nonzero arrays", func(t *testing.T) {
		a := [][][]int{{{3, 4}}}
		b := [][][]int{{{3, 4}}}
		cos, score := tpm_core.CosineSimWeights(1, k, n, a, b)
		if !floatsAlmostEqual(cos, 1.0, 1e-9) {
			t.Errorf("cosine = %v, want ~1.0", cos)
		}
		if score != 2 {
			t.Errorf("score = %v, want 2", score)
		}
	})

	t.Run("negated arrays", func(t *testing.T) {
		a := [][][]int{{{3, 4}}}
		b := [][][]int{{{-3, -4}}}
		cos, score := tpm_core.CosineSimWeights(1, k, n, a, b)
		if !floatsAlmostEqual(cos, -1.0, 1e-9) {
			t.Errorf("cosine = %v, want ~-1.0", cos)
		}
		if score != 0 {
			t.Errorf("score = %v, want 0", score)
		}
	})

	t.Run("orthogonal", func(t *testing.T) {
		a := [][][]int{{{1, 0}}}
		b := [][][]int{{{0, 1}}}
		cos, score := tpm_core.CosineSimWeights(1, k, n, a, b)
		if !floatsAlmostEqual(cos, 0, 1e-9) {
			t.Errorf("cosine = %v, want ~0", cos)
		}
		if score != 0 {
			t.Errorf("score = %v, want 0", score)
		}
	})

	t.Run("both all zero short-circuits norm but still scores", func(t *testing.T) {
		a := [][][]int{{{0, 0}}}
		b := [][][]int{{{0, 0}}}
		cos, score := tpm_core.CosineSimWeights(1, k, n, a, b)
		if cos != 0 {
			t.Errorf("cosine = %v, want 0", cos)
		}
		if score != 2 {
			t.Errorf("score = %v, want 2 (all positions coincidentally equal)", score)
		}
	})
}

func TestCreateRandomStimulusArray(t *testing.T) {
	t.Run("shape", func(t *testing.T) {
		stim := tpm_core.CreateRandomStimulusArray(4, 5, 3)
		if len(stim) != 4 {
			t.Fatalf("len(stim) = %v, want 4", len(stim))
		}
		for i, row := range stim {
			if len(row) != 5 {
				t.Errorf("len(stim[%d]) = %v, want 5", i, len(row))
			}
		}
	})

	t.Run("values in range and never zero", func(t *testing.T) {
		m := 4
		stim := tpm_core.CreateRandomStimulusArray(20, 20, m)
		for i, row := range stim {
			for j, v := range row {
				if v == 0 {
					t.Fatalf("stim[%d][%d] = 0, values must never be 0", i, j)
				}
				if v < -m || v > m || (v > -1 && v < 1) {
					t.Fatalf("stim[%d][%d] = %v, want in [-%d,-1] U [1,%d]", i, j, v, m, m)
				}
			}
		}
	})

	t.Run("m=1 only produces -1 or 1", func(t *testing.T) {
		stim := tpm_core.CreateRandomStimulusArray(10, 10, 1)
		for i, row := range stim {
			for j, v := range row {
				if v != -1 && v != 1 {
					t.Fatalf("stim[%d][%d] = %v, want -1 or 1", i, j, v)
				}
			}
		}
	})
}

func TestCreateRandomLayerWeightsArray(t *testing.T) {
	t.Run("shape", func(t *testing.T) {
		w := tpm_core.CreateRandomLayerWeightsArray(4, 5, 3)
		if len(w) != 4 {
			t.Fatalf("len(w) = %v, want 4", len(w))
		}
		for i, row := range w {
			if len(row) != 5 {
				t.Errorf("len(w[%d]) = %v, want 5", i, len(row))
			}
		}
	})

	t.Run("values in [-l,l] inclusive", func(t *testing.T) {
		l := 4
		w := tpm_core.CreateRandomLayerWeightsArray(20, 20, l)
		for i, row := range w {
			if !allInRange(row, -l, l) {
				t.Fatalf("w[%d] = %v, want all values in [-%d,%d]", i, row, l, l)
			}
		}
	})

	t.Run("l=0 forces every value to 0", func(t *testing.T) {
		w := tpm_core.CreateRandomLayerWeightsArray(5, 5, 0)
		for i, row := range w {
			for j, v := range row {
				if v != 0 {
					t.Fatalf("w[%d][%d] = %v, want 0 when l=0", i, j, v)
				}
			}
		}
	})
}

func TestGetNetworkDataSize(t *testing.T) {
	t.Run("hand-computed multi-layer sum", func(t *testing.T) {
		got := tpm_core.GetNetworkDataSize(2, []int{3, 2}, []int{5, 3})
		want := 3*5 + 2*3
		if got != want {
			t.Errorf("GetNetworkDataSize() = %v, want %v", got, want)
		}
	})

	t.Run("H=0", func(t *testing.T) {
		if got := tpm_core.GetNetworkDataSize(0, []int{3, 2}, []int{5, 3}); got != 0 {
			t.Errorf("GetNetworkDataSize() = %v, want 0", got)
		}
	})

	t.Run("K/N longer than H only counts first H layers", func(t *testing.T) {
		got := tpm_core.GetNetworkDataSize(1, []int{3, 2}, []int{5, 3})
		want := 3 * 5
		if got != want {
			t.Errorf("GetNetworkDataSize() = %v, want %v", got, want)
		}
	})
}
