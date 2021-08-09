package solve

// AllPairingSets returns all pairing sets for alphabet. for instance with limit=4, this would return something like:
// [[0, 1] [2, 3]]
// [[0, 2] [1, 3]]
// [[0, 3] [1, 2]]
func AllPairingSets(limit int) [][][2]int {
	alphabet := make([]int, limit)
	for i := range alphabet {
		alphabet[i] = i
	}

	return allPairingSetsForAlphabet(alphabet)
}
func allPairingSetsForAlphabet(alphabet []int) [][][2]int {
	if len(alphabet) == 2 {
		return [][][2]int{{{alphabet[0], alphabet[1]}}}
	}

	first := alphabet[0]
	rest := alphabet[1:]
	var pairingsSet [][][2]int
	for i, v := range rest {
		pair := [2]int{first, v}
		withoutV := make([]int, len(rest)-1)
		copy(withoutV[:i], rest[:i])
		copy(withoutV[i:], rest[i+1:])

		// recursive call
		otherPairingsSet := allPairingSetsForAlphabet(withoutV)
		for _, otherPairings := range otherPairingsSet {
			thisPairing := make([][2]int, len(otherPairings)+1)
			thisPairing[0] = pair
			copy(thisPairing[1:], otherPairings)
			pairingsSet = append(pairingsSet, thisPairing)
		}
	}
	return pairingsSet
}
