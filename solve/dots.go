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
		tilesToSolutions[i] = g.SolveDot(tile)
	}

	// now merge them all together
	for i := 1; i < len(dotTiles); i++ {
		mergedIter := mergeSolutionsIters(tilesToSolutions[i-1], tilesToSolutions[i])
		uniqueIter := filterUnique(mergedIter)
		tilesToSolutions[i] = uniqueIter
	}

	return tilesToSolutions[len(dotTiles)-1]
}

// SolveDot returns a channel of solutions for a given dot tile.
func (g GridSolver) SolveDot(t gs.Tile) <-chan gs.TileSet {
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

	return g.solveDotRecur(t.Coord, gs.NewTileCoordSet(), numDots-knownEnabledTiles.Len())
}

func (g GridSolver) solveDotRecur(
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
