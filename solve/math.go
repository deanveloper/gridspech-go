package solve

func sliceIntersectGoalSolutions(s1, s2 []goalSolution) []goalSolution {
	var out []goalSolution
	for i1 := 0; i1 < len(s1)-1; i1++ {
		for i2 := i1 + 1; i2 < len(s1); i2++ {
			if s1[i1].eq(s2[i2]) {
				out = append(out, s1[i1])
			}
		}
	}

	return out
}

// AllPairingSets returns all pairing sets for alphabet. for instance with alphabet [1, 2, 3, 4], this would return something like:
// [1, 2] [3, 4]
// [1, 3] [2, 4]
// [1, 4] [2, 3]
func AllPairingSets(alphabet []int) [][][2]int {
	if len(alphabet) == 0 {
		return nil
	}

	var setsOfPairs [][][2]int
	forEachPairing(alphabet, func(start, end int, remaining []int) {
		setsOfPairs = append(setsOfPairs, allPairingSetsStartingWith(start, end, remaining))
	})

	return setsOfPairs
}
func allPairingSetsStartingWith(i1, i2 int, remaining []int) [][2]int {
	all := [][2]int{{i1, i2}}

	forEachPairing(remaining, func(start, end int, remaining []int) {
		all = append(all, allPairingSetsStartingWith(start, end, remaining)...)
	})

	return all
}
func forEachPairing(alphabet []int, forEach func(start, end int, remaining []int)) {
	for i1 := 0; i1 < len(alphabet)-1; i1++ {
		for i2 := i1 + 1; i2 < len(alphabet); i2++ {
			start, end := alphabet[i1], alphabet[i2]

			var remaining []int
			for _, each := range alphabet {
				if each != start && each != end {
					remaining = append(remaining, each)
				}
			}

			forEach(start, end, remaining)
		}
	}
}
