package solve

import (
	"fmt"

	"github.com/deanveloper/gridspech-go"
	gs "github.com/deanveloper/gridspech-go"
)

// SolveDots will return a slice of solutions for all of the dot tiles in g.
func (g GridSolver) SolveDots() <-chan gs.TileSet {

	// get all dot-related tiles
	dotTiles := g.Grid.TilesWith(func(o gs.Tile) bool {
		return o.Data.Type == gridspech.TypeDot1 || o.Data.Type == gridspech.TypeDot2 || o.Data.Type == gridspech.TypeDot3
	}).Slice()

	tilesToSolutions := make([]<-chan gs.TileSet, len(dotTiles))
	for i, tile := range dotTiles {
		tilesToSolutions[i] = g.solveDots(tile)
	}

	// now merge them all together
	for i := 1; i < len(dotTiles); i++ {
		mergedIter := mergeSolutionsIters(tilesToSolutions[i-1], tilesToSolutions[i])
		uniqueIter := filterUnique(mergedIter)
		tilesToSolutions[i] = uniqueIter
	}

	return tilesToSolutions[len(dotTiles)-1]
}

func mergeSolutionsIters(sols1, sols2 <-chan gs.TileSet) <-chan gs.TileSet {
	iter := make(chan gs.TileSet, 50)

	go func() {
		// read sols2 into a slice
		var sols2slice []gs.TileSet
		for sol2 := range sols2 {
			sols2slice = append(sols2slice, sol2)
		}

		// merge
		for sol1 := range sols1 {

		nextSolution:
			for _, sol2 := range sols2slice {
				var merged gs.TileSet

				// do not merge if they have any tiles with unmatched colors
				for _, t1 := range sol1.Slice() {
					for _, t2 := range sol2.Slice() {
						if t1.Coord == t2.Coord && t1.Data.Color != t2.Data.Color {
							continue nextSolution
						}
					}
				}

				merged.Merge(sol1)
				merged.Merge(sol2)
				iter <- merged
			}
		}
		close(iter)
	}()

	return iter
}

func filterUnique(in <-chan gs.TileSet) <-chan gs.TileSet {
	filtered := make(chan gs.TileSet, 200)

	go func() {
		var alreadySeen []gs.TileSet
		for newSolution := range in {
			unique := true
			for _, seen := range alreadySeen {
				if newSolution.Eq(seen) {
					unique = false
					break
				}
			}
			if unique {
				alreadySeen = append(alreadySeen, newSolution)
				filtered <- newSolution
			}
		}
		close(filtered)
	}()

	return filtered
}

// there are very few valid solutions for an individual tile, so this just returns a slice
func (g GridSolver) solveDots(t gs.Tile) <-chan gs.TileSet {
	var numDots int

	switch t.Data.Type {
	case gs.TypeDot1:
		numDots = 1
	case gs.TypeDot2:
		numDots = 2
	case gs.TypeDot3:
		numDots = 3
	default:
		panic(fmt.Sprint("invalid type", t.Data.Type))
	}

	knownEnabledTiles := g.Grid.NeighborSetWith(t.Coord, func(o gs.Tile) bool {
		return o.Data.Color != gs.ColorNone && !g.UnknownTiles.Has(o.Coord)
	})

	return g.solveDotsRecur(t.Coord, gs.NewTileCoordSet(), numDots-knownEnabledTiles.Len())
}

func (g GridSolver) solveDotsRecur(
	t gs.TileCoord,
	tilesBeingUsed gs.TileCoordSet,
	numDots int,
) <-chan gs.TileSet {

	ch := make(chan gs.TileSet, 4)

	go func() {

		defer close(ch)

		if numDots < 0 {
			return
		}
		if numDots == 0 {
			ch <- gs.NewTileSet()
			return
		}

		unknownNeighbors := g.Grid.NeighborSliceWith(t, func(o gs.Tile) bool {
			return g.UnknownTiles.Has(o.Coord) && !tilesBeingUsed.Has(o.Coord)
		})

		// if there are not enough unknown neighbors to fulfil this dot, then there are no solutions
		if numDots > len(unknownNeighbors) {
			return
		}

		for perm := range Permutation(g.Grid.MaxColors, len(unknownNeighbors)) {
			var numNonZero int
			for _, i := range perm {
				if i > 0 {
					numNonZero++
				}
			}
			if numNonZero != numDots {
				continue
			}

			var result gs.TileSet
			for i, c := range perm {
				tCopy := unknownNeighbors[i]
				tCopy.Data.Color = gs.TileColor(c)
				result.Add(tCopy)
			}
			ch <- result
		}
	}()

	return ch
}
