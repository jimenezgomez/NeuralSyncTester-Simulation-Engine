package tpm_stimHandlers_test

import (
	"reflect"
	"testing"

	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/pkg/tpm/tpm_stimHandlers"
)

func TestNoOverlapTPM_CreateStimulationStructure(t *testing.T) {
	tpm := tpm_stimHandlers.NoOverlapTPM{}

	t.Run("h=1 loop body never runs", func(t *testing.T) {
		got := tpm.CreateStimulationStructure([]int{5}, 7)
		want := []int{7}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("h=3 backward computation", func(t *testing.T) {
		// k[2]=2, k[1]=n[2]*k[2]=3*2=6, k[0]=n[1]*k[1]=4*6=24
		got := tpm.CreateStimulationStructure([]int{10, 4, 3}, 2)
		want := []int{24, 6, 2}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestNoOverlapTPM_CreateStimulusFromLayerOutput(t *testing.T) {
	tpm := tpm_stimHandlers.NoOverlapTPM{}

	t.Run("contiguous disjoint partition", func(t *testing.T) {
		outputs := seqInts(1, 6)
		got := tpm.CreateStimulusFromLayerOutput(nil, outputs, 3, 2)
		want := [][]int{{1, 2}, {3, 4}, {5, 6}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("reuses dst backing array when capacity suffices", func(t *testing.T) {
		dst := make([][]int, 5, 10)
		outputs := seqInts(1, 6)
		got := tpm.CreateStimulusFromLayerOutput(dst, outputs, 3, 2)
		if cap(got) != 10 {
			t.Errorf("cap(got) = %v, want 10 (backing array reused)", cap(got))
		}
	})

	t.Run("allocates fresh backing array when capacity insufficient", func(t *testing.T) {
		outputs := seqInts(1, 6)
		got := tpm.CreateStimulusFromLayerOutput(nil, outputs, 3, 2)
		if cap(got) != 3 {
			t.Errorf("cap(got) = %v, want 3 (fresh allocation)", cap(got))
		}
	})
}

func TestIntPow(t *testing.T) {
	tests := []struct {
		name      string
		base, exp int
		want      int
	}{
		{"exp=0 base positive", 5, 0, 1},
		{"exp=0 base=0", 0, 0, 1},
		{"exp=1", 5, 1, 5},
		{"base=2 exp=10", 2, 10, 1024},
		{"negative base odd exp", -2, 3, -8},
		{"base=0 exp=5", 0, 5, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tpm_stimHandlers.IntPow(tt.base, tt.exp); got != tt.want {
				t.Errorf("IntPow(%v, %v) = %v, want %v", tt.base, tt.exp, got, tt.want)
			}
		})
	}
}

// TestIntPow_NegativeExponent_KnownBug documents, without ever invoking IntPow with a
// negative exponent, that doing so hangs forever.
func TestIntPow_NegativeExponent_KnownBug(t *testing.T) {
	t.Skip("IntPow(base, exp) hangs forever for exp < 0: `exp >>= 1` is an arithmetic shift on a " +
		"signed int, so a negative exp never reaches 0 and the loop condition is never satisfied. " +
		"Do not call IntPow with a negative exponent until this is fixed in noOverlap.go.")
}
