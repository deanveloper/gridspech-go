package solve

import (
	"fmt"
	"sync"

	"github.com/deanveloper/gridspech-go"
	gs "github.com/deanveloper/gridspech-go"
)

// Dots will return a channel of solutions for all of the dot tiles in g.
func Dots(g GridSolver, maxColors int) <-chan GridSolver {

	// get all dot-related tiles
	dotTiles := g.RawGrid.TilesWith(func(o gs.Tile) bool {
		return o.Type == gridspech.TypeDot1 || o.Type == gridspech.TypeDot2 || o.Type == gridspech.TypeDot3
	}).Slice()

	var wg sync.WaitGroup
	ch := make(chan GridSolver)

	wg.Add(len(dotTiles))
	go func() {
		wg.Wait()
		close(ch)
	}()

	// solve each dot related tile
	for _, dotTile := range dotTiles {
		dotTile := dotTile
		go func() {
			solutions := g.solveDots(dotTile, maxColors)
			for _, sol := range solutions {
				ch <- sol
			}
			wg.Done()
		}()
	}

	return ch
}

// there are very few valid solutions for an individual tile, so this just returns a slice
func (g GridSolver) solveDots(t gs.Tile, maxColors int) []GridSolver {
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

	return g.solveDotsRecur(t, maxColors, numDots-enabledTiles.Len())
}

func (g GridSolver) solveDotsRecur(t gs.Tile, maxColors int, remainingDots int) []GridSolver {

	// base case: exactly 0 remaining dots means this tile is now valid, so this grid is a solution
	if remainingDots == 0 {
		return []GridSolver{{g.RawGrid.Clone()}}
	}

	var grids []GridSolver

	unknownNeighbors := g.RawGrid.NeighborsWith(t, func(o gs.Tile) bool {
		return o.Color == ColorUnknown
	})

	// if there are not enough unknown neighbors to fulfil this dot, then there are no solutions
	if remainingDots > unknownNeighbors.Len() {
		return nil
	}

	// call recursively until dot is fulfilled
	for _, tile := range unknownNeighbors.Slice() {
		for c := 0; c < maxColors; c++ {
			newGrid := GridSolver{RawGrid: g.RawGrid.Clone()}
			newGrid.RawGrid.Tiles[tile.X][tile.Y].Color = gs.TileColor(c)
			moreGrids := newGrid.solveDotsRecur(t, maxColors, remainingDots-1)
			grids = append(grids, moreGrids...)
		}
	}

	return grids
}
