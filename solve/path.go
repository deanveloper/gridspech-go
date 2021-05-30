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
//   3. have the same state as start.
func (g Grid) SolveGoals(start, end Tile) <-chan TileSet {
	ch := make(chan TileSet)
	if end.Sticky && start.State != end.State {
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
func (g Grid) dfsDirectPaths(prev, end Tile, path TileSet, ch chan<- TileSet) {
	neighbors := g.Neighbors(prev)
	for _, next := range neighbors.Slice() {
		if path.Has(next) {
			continue
		}

		if next == end {
			path.Add(next)
			ch <- path
			return
		}

		// represents neighbors with the same state
		prevNeighbors := g.NeighborsWith(prev, func(t Tile) bool {
			return t.State == prev.State || path.Has(t)
		})

		// in diagrams: p is prev, n is next, x is same State, o is diff State

		// we prune:
		// ooo
		// xpn
		// oxo
		if prevNeighbors.Len() == 2 && !prevNeighbors.Has(next) {
			continue
		}
		// we prune:
		// ooo
		// xpn
		// ooo
		// where n is a sticky element with different State
		if prevNeighbors.Len() == 1 && !prevNeighbors.Has(next) && next.Sticky {
			continue
		}

		nextNeighbors := g.NeighborsWith(prev, func(t Tile) bool {
			return t.State == prev.State || path.Has(t)
		})

		// prune if next will def be invalid
		if nextNeighbors.Len() > 2 {
			continue
		}
		// we prune when we are next to a goal of the same type (and that neighbor is not prev)
		for _, neighbor := range nextNeighbors.Slice() {
			if neighbor.Type == gridspech.Goal && neighbor != prev {
				continue
			}
		}

		// setup for recursion
		path.Add(next)

		// RECURSION
		g.dfsDirectPaths(next, end, path, ch)

		// recursion takedown
		path.Remove(next)
	}
}
