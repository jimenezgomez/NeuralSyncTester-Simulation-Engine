package config_manager

// SolveSize_NoOverlap finds all (K, N) combinations of exactly H layers whose weight count equals target.
// Constraints: N[i] >= 2, Kh >= 1, Kh <= target/2.
// It works recursively: it calls the search funciton (Defined inside this func) for each layer
// --> Slowly building the N array on each depth pass, for each K from 1 to the max value
func SolveSize_Overlap(H, targetSize, minN0, maxN0 int, overlapType string) []Combination {
	var results []Combination

	// NRev accumulates N values in reverse order: NRev[0]=N[H-1], NRev[1]=N[H-2], ...
	var search func(depth int, K []int, remaining, Ncur, N0 int)
	search = func(depth int, K []int, remaining, Ncur, N0 int) {
		if depth == H {
			if remaining == 0 {
				var N []int
				switch overlapType {
				case "FULL OVERLAP":
					N = buildNFullOverlap(K, N0)
				case "PARTIAL OVERLAP":
					N = buildNPartialOverlap(K, N0)
				}
				results = append(results, Combination{K: K, N: N})
			}
			return
		}

		// Kcur is K for this layer; each valid Ni consumes Kcur*Ni from the budget.
		// The next layer to the left will have K = Kcur * Ni.
		maxK := remaining / Ncur
		for Ki := 2; Ki <= maxK; Ki++ {
			search(depth+1, append(K, Ki), remaining-Ncur*Ki, Ncur*Ki, N0)
		}
	}

	if minN0 == 0 {
		minN0 = 2
	}

	if maxN0 == 0 {
		maxN0 = targetSize / 2
	}

	// Kh is bounded by target/2 (since every layer contributes at least Kh*2).
	for N0 := minN0; N0 <= maxN0; N0++ {
		// Start filling from the last layer (depth=0 in NRev = layer H-1).
		// Kcur = Kh for the last layer.
		search(0, []int{}, targetSize, N0, N0)
	}

	return results
}
