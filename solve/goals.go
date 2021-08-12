package solve

import (
	"sync"

	gs "github.com/deanveloper/gridspech-go"
)

// Goals will return a channel of solutions for all the goal tiles in g
func Goals(g GridSolver, maxColors int) <-chan gs.TileSet {

	iter := make(chan gs.TileSet, 4)

	go func() {
		g.solveGoals(maxColors, iter)
		close(iter)
	}()

	return iter
}

func (g GridSolver) solveGoals(maxColors int, ch chan<- gs.TileSet) {
	goalTiles := g.RawGrid.TilesWith(func(o gs.Tile) bool {
		return o.Data.Type == gs.TypeGoal
	}).Slice()

	var pairsToSolutionMx sync.Mutex
	pairsToSolutions := make(map[[2]gs.Tile][]gs.TileSet)
	var wg sync.WaitGroup
	for i1 := 0; i1 < len(goalTiles)-1; i1++ {
		for i2 := i1 + 1; i2 < len(goalTiles); i2++ {
			goalPair := [2]gs.Tile{goalTiles[i1], goalTiles[i2]}
			wg.Add(1)
			go func() {
				for c := 0; c < maxColors; c++ {
					for path := range g.SolvePath(goalPair[0], goalPair[1], gs.TileColor(c)) {
						pairsToSolutionMx.Lock()
						pairsToSolutions[goalPair] = append(pairsToSolutions[goalPair], path)
						pairsToSolutionMx.Unlock()
					}
				}
				wg.Done()
			}()
		}
	}
	wg.Wait()

	// now we get solutions for each pairing
	allGoalPairings := allTilePairingSets(goalTiles)
	for _, pairing := range allGoalPairings {
		pairingSolutions := pairsToSolutions[pairing[0]]
		for pairIndex := 1; pairIndex < len(pairing); pairIndex++ {
			pair := pairing[pairIndex]
			var tilesToValidate []gs.Tile
			for i := 0; i <= pairIndex; i++ {
				tilesToValidate = append(tilesToValidate, pairing[i][0], pairing[i][1])
			}

			result := mergeSolutionsSlices(pairingSolutions, pairsToSolutions[pair])
			result = removeIfNonUnique(result)
			result = removeIfInvalid(g, tilesToValidate, result)
			pairingSolutions = result
		}
		for _, solution := range pairingSolutions {
			ch <- solution
		}
	}
}

func mergeSolutionsSlices(sols1, sols2 []gs.TileSet) []gs.TileSet {
	var result []gs.TileSet
	for _, sol1 := range sols1 {
		for _, sol2 := range sols2 {
			var merged gs.TileSet
			merged.Merge(sol1)
			merged.Merge(sol2)
			result = append(result, merged)
		}
	}
	return result
}

func removeIfInvalid(g GridSolver, tilesToValidate []gs.Tile, in []gs.TileSet) []gs.TileSet {
	var validSolutions []gs.TileSet

	base := g.Grid()
	for _, solution := range in {
		newBase := base.Clone()
		newBase.ApplyTileSet(solution)

		allValid := true
		for _, tile := range tilesToValidate {
			if !newBase.ValidTile(tile.Coord) {
				allValid = false
				break
			}
		}
		if allValid {
			validSolutions = append(validSolutions, solution)
		}
	}

	return validSolutions
}

func removeIfNonUnique(in []gs.TileSet) []gs.TileSet {
	var filtered []gs.TileSet

	for _, solution := range in {
		unique := true
		for _, seen := range filtered {
			if solution.Eq(seen) {
				unique = false
				break
			}
		}
		if unique {
			filtered = append(filtered, solution)
		}
	}

	return filtered
}

func allTilePairingSets(tiles []gs.Tile) [][][2]gs.Tile {

	pairingSets := AllPairingSets(len(tiles))
	tilePairingSets := make([][][2]gs.Tile, len(pairingSets))
	for i, pairing := range pairingSets {
		tilePairings := make([][2]gs.Tile, len(pairing))
		for p, pair := range pairing {
			tilePairings[p] = [2]gs.Tile{tiles[pair[0]], tiles[pair[1]]}
		}
		tilePairingSets[i] = tilePairings
	}
	return tilePairingSets
}
