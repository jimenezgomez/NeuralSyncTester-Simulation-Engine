package tpm_learnRules_test

import (
	"reflect"
	"testing"

	"github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/pkg/tpm/tpm_learnRules"
)

type namedRule struct {
	name string
	rule tpm_learnRules.TPMLearnRuleHandler
}

func allRules() []namedRule {
	return []namedRule{
		{"Hebbian", tpm_learnRules.HebbianLearnRule{}},
		{"AntiHebbian", tpm_learnRules.AntiHebbianLearnRule{}},
		{"RandomWalk", tpm_learnRules.RandomWalkLearnRule{}},
	}
}

func TestTPMLearnLayer_Disagreement_WeightsUnchanged(t *testing.T) {
	weights := [][]int{{1, -1}, {2, -2}}
	stim := [][]int{{3, 4}, {5, 6}}
	outputs := []int{1, 1}
	k, n, l := 2, 2, 10
	outputA, outputB := 1, -1 // agreement = Heaviside(-1) = 0

	for _, nr := range allRules() {
		t.Run(nr.name, func(t *testing.T) {
			before := cloneMatrix(weights)
			w := cloneMatrix(weights)
			nr.rule.TPMLearnLayer(k, n, l, w, stim, outputs, outputA, outputB)
			if !reflect.DeepEqual(w, before) {
				t.Errorf("weights changed on disagreement: got %v, want unchanged %v", w, before)
			}
		})
	}
}

func TestTPMLearnLayer_PartialNeuronMatch(t *testing.T) {
	weights := [][]int{{1, 2}, {3, 4}}
	stim := [][]int{{1, -1}, {2, -2}}
	outputs := []int{1, -1} // neuron 0 matches output_a, neuron 1 doesn't
	k, n, l := 2, 2, 10
	outputA, outputB := 1, 1 // agreement = 1

	wantRow1Unchanged := []int{3, 4}
	wantRow0 := map[string][]int{
		"Hebbian":     {2, 1},
		"AntiHebbian": {0, 3},
		"RandomWalk":  {2, 1},
	}

	for _, nr := range allRules() {
		t.Run(nr.name, func(t *testing.T) {
			w := cloneMatrix(weights)
			nr.rule.TPMLearnLayer(k, n, l, w, stim, outputs, outputA, outputB)
			if !reflect.DeepEqual(w[1], wantRow1Unchanged) {
				t.Errorf("non-matching neuron row changed: got %v, want %v", w[1], wantRow1Unchanged)
			}
			if !reflect.DeepEqual(w[0], wantRow0[nr.name]) {
				t.Errorf("matching neuron row = %v, want %v", w[0], wantRow0[nr.name])
			}
		})
	}
}

func TestTPMLearnLayer_FullAgreementFullMatch(t *testing.T) {
	weights := [][]int{{1, 2}, {3, 4}}
	stim := [][]int{{1, -1}, {2, -2}}
	outputs := []int{-1, -1}
	k, n, l := 2, 2, 10
	outputA, outputB := -1, -1 // agreement = 1, both neurons match (output_a=-1)

	want := map[string][][]int{
		// Hebbian: w + stim*output_a = w - stim (output_a=-1)
		"Hebbian": {{0, 3}, {1, 6}},
		// AntiHebbian: w - stim*output_a = w + stim
		"AntiHebbian": {{2, 1}, {5, 2}},
		// RandomWalk: w + stim, ignoring output_a entirely -> diverges from Hebbian here
		"RandomWalk": {{2, 1}, {5, 2}},
	}

	for _, nr := range allRules() {
		t.Run(nr.name, func(t *testing.T) {
			w := cloneMatrix(weights)
			nr.rule.TPMLearnLayer(k, n, l, w, stim, outputs, outputA, outputB)
			if !reflect.DeepEqual(w, want[nr.name]) {
				t.Errorf("got %v, want %v", w, want[nr.name])
			}
		})
	}

	// RandomWalk must diverge from Hebbian here since it ignores output_a's sign.
	hebbianW := cloneMatrix(weights)
	tpm_learnRules.HebbianLearnRule{}.TPMLearnLayer(k, n, l, hebbianW, stim, outputs, outputA, outputB)
	randomWalkW := cloneMatrix(weights)
	tpm_learnRules.RandomWalkLearnRule{}.TPMLearnLayer(k, n, l, randomWalkW, stim, outputs, outputA, outputB)
	if reflect.DeepEqual(hebbianW, randomWalkW) {
		t.Error("RandomWalk should diverge from Hebbian when output_a=-1, but results matched")
	}
}

func TestTPMLearnLayer_ClippingBoundary(t *testing.T) {
	l := 5

	t.Run("Hebbian upper bound clamps to +l", func(t *testing.T) {
		weights := [][]int{{5}}
		stim := [][]int{{1}}
		outputs := []int{1}
		tpm_learnRules.HebbianLearnRule{}.TPMLearnLayer(1, 1, l, weights, stim, outputs, 1, 1)
		if weights[0][0] != l {
			t.Errorf("weight = %v, want clamped to %v", weights[0][0], l)
		}
	})

	t.Run("AntiHebbian lower bound clamps to -l", func(t *testing.T) {
		weights := [][]int{{-5}}
		stim := [][]int{{1}}
		outputs := []int{1}
		tpm_learnRules.AntiHebbianLearnRule{}.TPMLearnLayer(1, 1, l, weights, stim, outputs, 1, 1)
		if weights[0][0] != -l {
			t.Errorf("weight = %v, want clamped to %v", weights[0][0], -l)
		}
	})
}

func TestTPMLearnLayer_MultiWeightIntegration(t *testing.T) {
	weights := [][]int{{1, -1}, {2, -2}}
	stim := [][]int{{3, -4}, {-5, 6}}
	outputs := []int{1, 1}
	k, n, l := 2, 2, 10
	outputA, outputB := 1, 1 // agreement=1, both neurons match

	want := [][]int{{4, -5}, {-3, 4}} // Hebbian: w + stim (output_a=1, matches=1, agreement=1)

	w := cloneMatrix(weights)
	tpm_learnRules.HebbianLearnRule{}.TPMLearnLayer(k, n, l, w, stim, outputs, outputA, outputB)
	if !reflect.DeepEqual(w, want) {
		t.Errorf("got %v, want %v", w, want)
	}
}
