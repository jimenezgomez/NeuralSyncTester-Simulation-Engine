package tpm_stimHandlers_test

import (
	"reflect"
	"testing"

	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/pkg/tpm/tpm_stimHandlers"
)

func TestFullOverlapTPM_CreateStimulationStructure(t *testing.T) {
	tpm := tpm_stimHandlers.FullOverlapTPM{}

	t.Run("h=1", func(t *testing.T) {
		got := tpm.CreateStimulationStructure([]int{3}, 10)
		want := []int{10}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("h=3", func(t *testing.T) {
		got := tpm.CreateStimulationStructure([]int{3, 4, 2}, 1000)
		want := []int{1000, 3, 4}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestFullOverlapTPM_CreateStimulusFromLayerOutput(t *testing.T) {
	tpm := tpm_stimHandlers.FullOverlapTPM{}

	t.Run("every row is an identical broadcast copy", func(t *testing.T) {
		outputs := seqInts(1, 3)
		got := tpm.CreateStimulusFromLayerOutput(nil, outputs, 3, 3)
		want := [][]int{{1, 2, 3}, {1, 2, 3}, {1, 2, 3}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("reuses dst backing array when capacity suffices", func(t *testing.T) {
		dst := make([][]int, 5, 10)
		outputs := seqInts(1, 2)
		got := tpm.CreateStimulusFromLayerOutput(dst, outputs, 3, 2)
		if cap(got) != 10 {
			t.Errorf("cap(got) = %v, want 10 (backing array reused)", cap(got))
		}
	})

	t.Run("allocates fresh backing array when capacity insufficient", func(t *testing.T) {
		outputs := seqInts(1, 2)
		got := tpm.CreateStimulusFromLayerOutput(nil, outputs, 3, 2)
		if cap(got) != 3 {
			t.Errorf("cap(got) = %v, want 3 (fresh allocation)", cap(got))
		}
	})

	t.Run("k_h=0 produces empty result", func(t *testing.T) {
		outputs := seqInts(1, 2)
		got := tpm.CreateStimulusFromLayerOutput(nil, outputs, 0, 2)
		if len(got) != 0 {
			t.Errorf("len(got) = %v, want 0", len(got))
		}
	})
}
