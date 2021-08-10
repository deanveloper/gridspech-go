package solve

import (
	"fmt"

	"github.com/deanveloper/gridspech-go"
	gs "github.com/deanveloper/gridspech-go"
)

// Dots will return a slice of solutions for all of the dot tiles in g.
func Dots(g GridSolver, maxColors int) []gs.TileSet {

	// get all dot-related tiles
	dotTiles := g.RawGrid.TilesWith(func(o gs.Tile) bool {
		return o.Type == gridspech.TypeDot1 || o.Type == gridspech.TypeDot2 || o.Type == gridspech.TypeDot3
	}).Slice()

	runningSolutions := []gs.TileSet{gs.NewTileSet()}

	// solve each dot related tile
	for currentIndex, currentTile := range dotTiles {
		newSolutions := g.solveDots(currentTile, maxColors)
		mergedSolutions := mergeSolutions(g, runningSolutions, newSolutions)
		runningSolutions = nil

		// check validity of each new solution
		for _, solution := range mergedSolutions {

			var anyDotsInvalid bool

			newGrid := g.Grid()
			for _, tile := range solution.Slice() {
				newGrid.Tiles[tile.X][tile.Y].Color = tile.Color
			}
			for prevIndex := 0; prevIndex < currentIndex; prevIndex++ {
				prevTile := dotTiles[prevIndex]
				if !newGrid.ValidTile(prevTile) {
					anyDotsInvalid = false
					break
				}
			}

			if !anyDotsInvalid {
				runningSolutions = append(runningSolutions, solution)
			}
		}
	}

	return runningSolutions
}

func mergeSolutions(g GridSolver, prevSols, newSols []gs.TileSet) []gs.TileSet {
	var solutions []gs.TileSet
	for _, sol1 := range prevSols {
		for _, sol2 := range newSols {

			// merge the tilesets together
			var merged gs.TileSet
			merged.Merge(sol1)
			merged.Merge(sol2)
			solutions = append(solutions, merged)
		}
	}
	return solutions
}

// there are very few valid solutions for an individual tile, so this just returns a slice
func (g GridSolver) solveDots(t gs.Tile, maxColors int) []gs.TileSet {
	var numDots int

	switch t.Type {
	case gs.TypeDot1:
		numDots = 1
	case gs.TypeDot2:
		numDots = 2
	case gs.TypeDot3:
		numDots = 3
	default:
		panic(fmt.Sprint("invalid type", t.Type))
	}

	enabledTiles := g.RawGrid.NeighborsWith(t, func(o gs.Tile) bool {
		return o.Color != ColorUnknown && o.Color != gs.ColorNone
	})

	return solveDotsRecur(g.Clone(), t, maxColors, gs.NewTileSet(), numDots-enabledTiles.Len())
}

func solveDotsRecur(
	g GridSolver,
	t gs.Tile,
	maxColors int,
	runningSolution gs.TileSet,
	remainingDots int,
) []gs.TileSet {

	// base case: exactly 0 remaining dots means this tile is now valid, so the solution we have is the solution
	if remainingDots == 0 {
		var finalSolution gs.TileSet
		finalSolution.Merge(runningSolution)
		return []gs.TileSet{finalSolution}
	}

	unknownNeighbors := g.RawGrid.NeighborsWith(t, func(o gs.Tile) bool {
		return o.Color == ColorUnknown && !runningSolution.Has(o)
	})

	// if there are not enough unknown neighbors to fulfil this dot, then there are no solutions
	if remainingDots > unknownNeighbors.Len() {
		return nil
	}

	var solutions []gs.TileSet

	// call recursively until dot is fulfilled
	for _, tile := range unknownNeighbors.Slice() {

		// c=1 to avoid ColorNone
		for c := 1; c < maxColors; c++ {

			newTile := tile
			newTile.Color = gs.TileColor(c)

			runningSolution.Add(newTile)
			moreSolutions := solveDotsRecur(g, t, maxColors, runningSolution, remainingDots-1)
			runningSolution.Remove(newTile)

			solutions = append(solutions, moreSolutions...)
		}
	}

	fmt.Println(solutions)
	return solutions
}
