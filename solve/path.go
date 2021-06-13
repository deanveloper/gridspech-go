package solve

import gridspech "github.com/deanveloper/gridspech-go"

// Grid is an "extension" of gridspech.Grid with solving capabilities
type Grid struct {
	gridspech.Grid
}

// Tile is an alias for gridspech.Tile
type Tile = gridspech.Tile

// TileSet is an alias for gridspech.TileSet
type TileSet = gridspech.TileSet

// SolveGoals returns an channel of DFS direct paths from start to end.
// These paths will:
//   1. never contain a goal tile that isn't start or end.
//   2. never make a path that would cause start or end to become invalid Goal tiles.
//   3. have the same Color as start.
func (g Grid) SolveGoals(start, end Tile) <-chan TileSet {
	ch := make(chan TileSet)
	if end.Sticky && start.Color != end.Color {
		close(ch)
		return ch
	}
	go func() {
		var ts TileSet
		ts.Init()
		g.dfsDirectPaths(start, end, ts, ch)
		close(ch)
	}()
	return ch
}

// we do not iterate in any particular order since it does not matter.
// this function will only create direct paths, aka ones which would satisfy
// a Goal tile.
func (g Grid) dfsDirectPaths(prev, end Tile, path TileSet, ch chan<- TileSet) {
	neighbors := g.Neighbors(prev)
	for _, next := range neighbors.Slice() {

		// no circular paths
		if path.Has(next) {
			continue
		}

		// we cannot traverse into a Goal tile
		if next.Type == gridspech.TypeGoal {
			continue
		}

		if next == end {
			path.Add(next)
			ch <- path
			return
		}

		// represents neighbors with the same Color (or prospective Color)
		prevNeighbors := g.NeighborsWith(prev, func(t Tile) bool {
			return t.Color == prev.Color || path.Has(t)
		})

		// in diagrams: p is prev, n is next, x is same Color, o is diff Color

		// we prune:
		// ooo
		// xpn
		// oxo
		// (aka we will not create a new tile of a different color if
		// we already have 2 neighbors of the same color)
		if prevNeighbors.Len() == 2 && !prevNeighbors.Has(next) {
			continue
		}
		// we prune:
		// ooo
		// xpn
		// ooo
		// where n is a sticky element with different Color
		// (aka we cannot change a tile that is sticky)
		if prevNeighbors.Len() == 1 && !prevNeighbors.Has(next) && next.Sticky {
			continue
		}

		// setup for recursion
		path.Add(next)

		// RECURSION
		g.dfsDirectPaths(next, end, path, ch)

		// recursion takedown
		path.Remove(next)
	}
}
