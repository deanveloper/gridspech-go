package solve

import (
	gs "github.com/deanveloper/gridspech-go"
)

// ColorUnknown is a special color which represents
const ColorUnknown gs.TileColor = 'U' - 'A' - 1

// GridSolver represents a gs.Grid, but with a special "unknown" tile color.
type GridSolver struct {
	RawGrid gs.Grid
}

// NewGridSolver creates a GridSolver
func NewGridSolver(solving gs.Grid) GridSolver {
	newSolving := solving.Clone()
	for x := range newSolving.Tiles {
		for y := range newSolving.Tiles[x] {
			if !newSolving.Tiles[x][y].Sticky {
				newSolving.Tiles[x][y].Color = ColorUnknown
			}
		}
	}
	return GridSolver{RawGrid: newSolving}
}

// Grid returns the underlying gridspech.Grid, with unknown tiles replaced with gridspech.ColorNone
func (g GridSolver) Grid() gs.Grid {
	newSolving := g.RawGrid.Clone()
	for x := range newSolving.Tiles {
		for y := range newSolving.Tiles[x] {
			if newSolving.Tiles[x][y].Color == ColorUnknown {
				newSolving.Tiles[x][y].Color = gs.ColorNone
			}
		}
	}
	return newSolving
}

// Clone creates a copy of the underlying gridspech.Grid
func (g GridSolver) Clone() GridSolver {
	return GridSolver{RawGrid: g.RawGrid.Clone()}
}

// SolvePath returns an channel of direct paths from start to end.
// These paths will:
//   1. never contain a goal tile that isn't start or end.
//   2. never make a path that would cause start or end to become invalid Goal tiles.
//   3. have the same Color as start.
func (g GridSolver) SolvePath(start, end gs.Tile, color gs.TileColor) <-chan gs.TileSet {
	ch := make(chan gs.TileSet)
	if start.Sticky && color != start.Color {
		close(ch)
		return ch
	}
	if end.Sticky && color != end.Color {
		close(ch)
		return ch
	}
	go func() {
		var ts gs.TileSet
		ts.Add(start)
		g.dfsDirectPaths(color, start, end, ts, ch)
		close(ch)
	}()
	return ch
}

// we do not iterate in any particular order since it does not matter.
// this function will only create direct paths, aka ones which would satisfy
// a Goal tile.
func (g GridSolver) dfsDirectPaths(color gs.TileColor, prev, end gs.Tile, path gs.TileSet, ch chan<- gs.TileSet) {

	// possible next tiles include untraversed tiles, and tiles of the same color
	possibleNext := g.RawGrid.NeighborsWith(prev, func(o gs.Tile) bool {
		return o.Color == ColorUnknown || o.Color == color
	})

	for _, next := range possibleNext.Slice() {
		// no circular paths
		if path.Has(next) {
			continue
		}

		// prev's neighbors with same color, including `next`
		prevNeighborsSameColor := g.RawGrid.NeighborsWith(prev, func(o gs.Tile) bool {
			return o.Color == color || path.Has(o) || o == next
		})
		// make sure that we never have an invalid path, ever
		if prevNeighborsSameColor.Len() > 2 {
			continue
		}
		if prev.Type == gs.TypeGoal && prevNeighborsSameColor.Len() > 1 {
			continue
		}

		// we cannot traverse through a Goal tile that is not the end tile
		if next.Type == gs.TypeGoal && next != end {
			continue
		}

		// we found a possible solution
		if next == end {

			// make sure the goal only has 1 neighbor of the same color
			if end.Type == gs.TypeGoal {
				endNeighbors := g.RawGrid.NeighborsWith(end, func(o gs.Tile) bool {
					return o.Color == color || path.Has(o)
				})
				if endNeighbors.Len() > 1 {
					continue
				}
			}

			var finalPath gs.TileSet
			finalPath.Merge(path)
			finalPath.Add(next)
			ch <- finalPath
			continue
		}

		var nextPath gs.TileSet
		nextPath.Merge(path)
		nextPath.Add(next)

		// RECURSION
		g.dfsDirectPaths(color, next, end, nextPath, ch)
	}
}
