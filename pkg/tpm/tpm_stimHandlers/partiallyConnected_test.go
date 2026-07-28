package tpm_stimHandlers_test

import (
	"reflect"
	"testing"

	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/pkg/tpm/tpm_stimHandlers"
)

func TestPartialOverlapTPM_CreateStimulationStructure(t *testing.T) {
	tpm := tpm_stimHandlers.PartialOverlapTPM{}

	t.Run("two layers strictly decreasing k", func(t *testing.T) {
		got := tpm.CreateStimulationStructure([]int{5, 3}, 10)
		want := []int{10, 3}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("three layers strictly decreasing k", func(t *testing.T) {
		got := tpm.CreateStimulationStructure([]int{5, 3, 2}, 10)
		want := []int{10, 3, 2}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("non-decreasing k returns nil", func(t *testing.T) {
		got := tpm.CreateStimulationStructure([]int{3, 3}, 10)
		if got != nil {
			t.Errorf("got %v, want nil", got)
		}
	})

	t.Run("increasing k returns nil", func(t *testing.T) {
		got := tpm.CreateStimulationStructure([]int{3, 5}, 10)
		if got != nil {
			t.Errorf("got %v, want nil", got)
		}
	})

	t.Run("h=1 loop does not execute", func(t *testing.T) {
		got := tpm.CreateStimulationStructure([]int{5}, 10)
		want := []int{10}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestPartialOverlapTPM_CreateStimulusFromLayerOutput(t *testing.T) {
	tpm := tpm_stimHandlers.PartialOverlapTPM{}

	t.Run("sliding window", func(t *testing.T) {
		outputs := seqInts(10, 5) // [10,11,12,13,14]
		got := tpm.CreateStimulusFromLayerOutput(nil, outputs, 3, 3)
		want := [][]int{{10, 11, 12}, {11, 12, 13}, {12, 13, 14}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("reuses dst backing array when capacity suffices", func(t *testing.T) {
		dst := make([][]int, 5, 10)
		outputs := seqInts(10, 5)
		got := tpm.CreateStimulusFromLayerOutput(dst, outputs, 3, 3)
		if cap(got) != 10 {
			t.Errorf("cap(got) = %v, want 10 (backing array reused)", cap(got))
		}
	})

	t.Run("allocates fresh backing array when capacity insufficient", func(t *testing.T) {
		outputs := seqInts(10, 5)
		got := tpm.CreateStimulusFromLayerOutput(nil, outputs, 3, 3)
		if cap(got) != 3 {
			t.Errorf("cap(got) = %v, want 3 (fresh allocation)", cap(got))
		}
	})
}
