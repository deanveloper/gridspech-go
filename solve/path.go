package solve

import (
	gs "github.com/deanveloper/gridspech-go"
)

// Grid is an "extension" of gridspech.Grid with solving capabilities
type Grid struct {
	gs.Grid
}

// Path returns an channel of DFS direct paths from start to end.
// These paths will:
//   1. never contain a goal tile that isn't start or end.
//   2. never make a path that would cause start or end to become invalid Goal tiles.
//   3. have the same Color as start.
func (g Grid) Path(start, end gs.Tile, color gs.TileColor) <-chan gs.TileSet {
	ch := make(chan gs.TileSet)
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
func (g Grid) dfsDirectPaths(color gs.TileColor, prev, end gs.Tile, path gs.TileSet, ch chan<- gs.TileSet) {
	neighbors := g.Neighbors(prev)

	for _, next := range neighbors.Slice() {

		// no circular paths
		if path.Has(next) {
			continue
		}

		// represents neighbors with the same Color (or prospective Color)
		prevNeighborsSameColor := g.NeighborsWith(prev, func(o gs.Tile) bool {
			return o.Color == color || path.Has(o) || o == next
		})

		// first make sure that we never have an invalid path, ever
		if prevNeighborsSameColor.Len() > 2 {
			continue
		}

		// cannot traverse into sticky tile of different color
		if next.Color != color && next.Sticky {
			continue
		}

		// we cannot traverse through a Goal tile that is not the end tile
		if next.Type == gs.TypeGoal && next != end {
			continue
		}

		// we found a possible solution
		if next == end {
			// goals can only have 1 neighbor of the same color
			endNeighbors := g.NeighborsWith(end, func(o gs.Tile) bool {
				return o.Color == color || path.Has(o)
			})
			if endNeighbors.Len() > 1 {
				continue
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
