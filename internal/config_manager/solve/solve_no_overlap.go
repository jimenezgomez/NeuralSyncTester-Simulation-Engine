package config_manager

// SolveSize_NoOverlap finds all (K, N) combinations of exactly H layers whose weight count equals target.
// Constraints: N[i] >= 2, Kh >= 1, Kh <= target/2.
// It works recursively: it calls the search funciton (Defined inside this func) for each layer
// --> Slowly building the N array on each depth pass, for each K from 1 to the max value
func SolveSize_NoOverlap(H, targetSize, minKh, maxKh int) []Combination {
	var results []Combination

	// NRev accumulates N values in reverse order: NRev[0]=N[H-1], NRev[1]=N[H-2], ...
	var search func(depth int, NRev []int, remaining, Kcur, Kh int)
	search = func(depth int, NRev []int, remaining, Kcur, Kh int) {
		if depth == H {
			if remaining == 0 {
				N := make([]int, H)
				for j, v := range NRev {
					N[H-1-j] = v
				}
				K := buildK(N, Kh)
				results = append(results, Combination{K: K, N: N})
			}
			return
		}

		// Kcur is K for this layer; each valid Ni consumes Kcur*Ni from the budget.
		// The next layer to the left will have K = Kcur * Ni.
		maxN := remaining / Kcur
		for Ni := 2; Ni <= maxN; Ni++ {
			search(depth+1, append(NRev, Ni), remaining-Kcur*Ni, Kcur*Ni, Kh)
		}
	}

	if minKh == 0 {
		minKh = 1
	}

	if maxKh == 0 {
		maxKh = targetSize / 2
	}

	// Kh is bounded by target/2 (since every layer contributes at least Kh*2).
	for Kh := minKh; Kh <= maxKh; Kh++ {
		// Start filling from the last layer (depth=0 in NRev = layer H-1).
		// Kcur = Kh for the last layer.
		search(0, []int{}, targetSize, Kh, Kh)
	}

	return results
}
