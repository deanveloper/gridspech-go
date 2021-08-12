package solve

import (
	"fmt"

	"github.com/deanveloper/gridspech-go"
	gs "github.com/deanveloper/gridspech-go"
)

// Dots will return a slice of solutions for all of the dot tiles in g.
func Dots(g GridSolver) <-chan gs.TileSet {

	// get all dot-related tiles
	dotTiles := g.RawGrid.TilesWith(func(o gs.Tile) bool {
		return o.Data.Type == gridspech.TypeDot1 || o.Data.Type == gridspech.TypeDot2 || o.Data.Type == gridspech.TypeDot3
	}).Slice()

	tilesToSolutions := make([]<-chan gs.TileSet, len(dotTiles))
	for i, tile := range dotTiles {
		tilesToSolutions[i] = g.solveDots(tile)
	}

	// now merge them all together
	for i := 0; i < len(dotTiles)-1; i++ {
		mergedIter := mergeSolutionsIters(tilesToSolutions[i], tilesToSolutions[i+1])
		uniqueIter := filterUnique(mergedIter)
		validIter := filterValid(g, dotTiles[:i+2], uniqueIter)
		tilesToSolutions[i+1] = validIter
	}

	return tilesToSolutions[len(dotTiles)-1]
}

func mergeSolutionsIters(sols1, sols2 <-chan gs.TileSet) <-chan gs.TileSet {
	iter := make(chan gs.TileSet, 4)

	go func() {
		// read sols2 into a slice
		var sols2slice []gs.TileSet
		for sol2 := range sols2 {
			sols2slice = append(sols2slice, sol2)
		}

		// merge
		for sol1 := range sols1 {
			for _, sol2 := range sols2slice {
				var merged gs.TileSet
				merged.Merge(sol1)
				merged.Merge(sol2)
				iter <- merged
			}
		}
		close(iter)
	}()

	return iter
}

func filterValid(g GridSolver, tilesToValidate []gs.Tile, sols <-chan gs.TileSet) <-chan gs.TileSet {
	filtered := make(chan gs.TileSet)

	go func() {
		base := g.Grid()
		for solution := range sols {
			newBase := base.Clone()
			newBase.ApplyTileSet(solution)

			allValid := true
			for _, dotTile := range tilesToValidate {
				if !newBase.ValidTile(dotTile.Coord) {
					allValid = false
					break
				}
			}
			if allValid {
				filtered <- solution
			}
		}
		close(filtered)
	}()

	return filtered
}

func filterUnique(in <-chan gs.TileSet) <-chan gs.TileSet {
	filtered := make(chan gs.TileSet, 4)

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

	enabledTiles := g.RawGrid.NeighborsWith(t.Coord, func(o gs.Tile) bool {
		return o.Data.Color != ColorUnknown && o.Data.Color != gs.ColorNone
	})

	return solveDotsRecur(g.Clone(), t, gs.NewTileCoordSet(), numDots-enabledTiles.Len())
}

func solveDotsRecur(
	g GridSolver,
	t gs.Tile,
	tilesBeingUsed gs.TileCoordSet,
	remainingDots int,
) <-chan gs.TileSet {
	ch := make(chan gs.TileSet, 4)

	go func() {

		defer close(ch)

		if remainingDots == 0 {
			ch <- gs.NewTileSet()
			return
		}

		// if there are not enough unknown neighbors to fulfil this dot, then there are no solutions
		unknownNeighbors := g.RawGrid.NeighborsWith(t.Coord, func(o gs.Tile) bool {
			return o.Data.Color == ColorUnknown && !tilesBeingUsed.Has(o.Coord)
		})
		if remainingDots > unknownNeighbors.Len() {
			return
		}

		for tile := range unknownNeighbors.Iter() {
			tilesBeingUsed.Add(tile.Coord)
			for subSolution := range solveDotsRecur(g, t, tilesBeingUsed, remainingDots-1) {
				// c=1 to avoid ColorNone
				for c := 1; c < g.RawGrid.MaxColors; c++ {

					newTile := tile
					newTile.Data.Color = gs.TileColor(c)

					var newSolution gs.TileSet
					newSolution.Merge(subSolution)
					newSolution.Add(newTile)
					ch <- newSolution
				}
			}
			tilesBeingUsed.Remove(tile.Coord)
		}
	}()

	return ch
}
