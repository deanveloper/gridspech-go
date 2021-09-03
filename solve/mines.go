package solve

import (
	"fmt"

	"github.com/deanveloper/gridspech-go"
	gs "github.com/deanveloper/gridspech-go"
)

// SolveMines will return a channel of solutions for all of the mine tiles in g.
func (g GridSolver) SolveMines() <-chan gs.TileSet {

	// get all mine-related tiles
	mineTiles := g.Grid.TilesWith(func(o gs.Tile) bool {
		return o.Data.Type == gridspech.TypeMine1 || o.Data.Type == gridspech.TypeMine2 || o.Data.Type == gridspech.TypeMine3
	}).Slice()

	tilesToSolutions := make([]<-chan gs.TileSet, len(mineTiles))
	for i, tile := range mineTiles {
		tilesToSolutions[i] = solveMines(g, tile)
	}

	// now merge them all together
	for i := 1; i < len(mineTiles); i++ {
		mergedIter := mergeSolutionsIters(tilesToSolutions[i-1], tilesToSolutions[i])
		uniqueIter := filterUnique(mergedIter)
		validIter := filterValidSoFar(g, mineTiles[:i+1], mineTiles[i], uniqueIter)
		tilesToSolutions[i] = validIter
	}

	return tilesToSolutions[len(mineTiles)-1]
}

func mergeSolutionsIters(sols1, sols2 <-chan gs.TileSet) <-chan gs.TileSet {
	iter := make(chan gs.TileSet, 200)

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

func filterValidSoFar(
	g GridSolver,
	previousTiles []gs.Tile,
	current gs.Tile,
	sols <-chan gs.TileSet,
) <-chan gs.TileSet {
	filtered := make(chan gs.TileSet, 200)

	go func() {
		defer close(filtered)
		for solution := range sols {
			newBase := g.Grid.Clone()
			newBase.ApplyTileSet(solution)

			var nearbyPreviousTiles []gs.Tile
			for _, mineTile := range previousTiles {
				xDist := mineTile.Coord.X - current.Coord.X
				yDist := mineTile.Coord.Y - current.Coord.Y
				if xDist < 0 {
					xDist = -xDist
				}
				if yDist < 0 {
					yDist = -yDist
				}
				if xDist+yDist <= 2 {
					nearbyPreviousTiles = append(nearbyPreviousTiles, mineTile)
				}
			}

			allValid := true
			for _, mineTile := range nearbyPreviousTiles {
				if !newBase.ValidTile(mineTile.Coord) {
					allValid = false
					break
				}
			}
			if allValid {
				filtered <- solution
			}
		}
	}()

	return filtered
}

// there are very few valid solutions for an individual tile, so this just returns a slice
func solveMines(g GridSolver, t gs.Tile) <-chan gs.TileSet {
	var mineNum int

	switch t.Data.Type {
	case gs.TypeMine1:
		mineNum = 1
	case gs.TypeMine2:
		mineNum = 2
	case gs.TypeMine3:
		mineNum = 3
	default:
		panic(fmt.Sprint("invalid type", t.Data.Type))
	}

	enabledTiles := g.Grid.NeighborsWith(t.Coord, func(o gs.Tile) bool {
		return o.Data.Color != gs.ColorNone
	})

	return solveMinesRecur(g, t, gs.NewTileCoordSet(), mineNum-enabledTiles.Len())
}

func solveMinesRecur(
	g GridSolver,
	t gs.Tile,
	tilesBeingUsed gs.TileCoordSet,
	remainingTiles int,
) <-chan gs.TileSet {
	ch := make(chan gs.TileSet, 4)

	go func() {

		defer close(ch)

		if remainingTiles == 0 {
			ch <- gs.NewTileSet()
			return
		}

		// if there are not enough unknown neighbors to fulfil this mine, then there are no solutions
		unknownNeighbors := g.Grid.NeighborsWith(t.Coord, func(o gs.Tile) bool {
			return g.UnknownTiles.Has(o.Coord) && !tilesBeingUsed.Has(o.Coord)
		})
		if remainingTiles > unknownNeighbors.Len() {
			return
		}

		for tile := range unknownNeighbors.Iter() {
			tilesBeingUsed.Add(tile.Coord)
			for subSolution := range solveMinesRecur(g, t, tilesBeingUsed, remainingTiles-1) {
				// c=1 to avoid ColorNone
				for c := 1; c < g.Grid.MaxColors; c++ {

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
