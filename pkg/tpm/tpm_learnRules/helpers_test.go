package tpm_learnRules_test

import "github.com/jimenezgomez/NeuralSyncTester-Simulation-Engine/pkg/tpm/tpm_learnRules"

func cloneMatrix(m [][]int) [][]int {
	out := make([][]int, len(m))
	for i, row := range m {
		out[i] = append([]int(nil), row...)
	}
	return out
}

var (
	_ tpm_learnRules.TPMLearnRuleHandler = tpm_learnRules.HebbianLearnRule{}
	_ tpm_learnRules.TPMLearnRuleHandler = tpm_learnRules.AntiHebbianLearnRule{}
	_ tpm_learnRules.TPMLearnRuleHandler = tpm_learnRules.RandomWalkLearnRule{}
)
