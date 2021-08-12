package solve

import (
	gs "github.com/deanveloper/gridspech-go"
)

// ColorUnknown is a special color which represents a tile whose final color is not known.
const ColorUnknown gs.TileColor = 100

// GridSolver represents a gs.Grid, but with a special "unknown" tile color.
type GridSolver struct {
	RawGrid gs.Grid
}

// NewGridSolver creates a GridSolver
func NewGridSolver(solving gs.Grid) GridSolver {
	newSolving := solving.Clone()
	nonSticky := newSolving.TilesWith(func(o gs.Tile) bool {
		return !o.Data.Sticky
	})
	for tile := range nonSticky.Iter() {
		newSolving.TileAtCoord(tile.Coord).Data.Color = ColorUnknown
	}
	return GridSolver{RawGrid: newSolving}
}

// Grid returns the underlying gridspech.Grid, with unknown tiles replaced with gridspech.ColorNone
func (g GridSolver) Grid() gs.Grid {
	newSolving := g.RawGrid.Clone()
	for x := range newSolving.Tiles {
		for y := range newSolving.Tiles[x] {
			if newSolving.TileAt(x, y).Data.Color == ColorUnknown {
				newSolving.Tiles[x][y].Data.Color = gs.ColorNone
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
func (g GridSolver) SolvePath(start, end gs.TileCoord, color gs.TileColor) <-chan gs.TileSet {
	tileCoordSetIter := make(chan gs.TileCoordSet)
	go func() {
		defer close(tileCoordSetIter)
		startTile, endTile := *g.RawGrid.TileAtCoord(start), *g.RawGrid.TileAtCoord(end)
		if startTile.Data.Sticky && color != startTile.Data.Color {
			return
		}
		if endTile.Data.Sticky && color != endTile.Data.Color {
			return
		}
		g.dfsDirectPaths(color, startTile, endTile, gs.NewTileCoordSet(start), tileCoordSetIter)
	}()

	onGridIter := tileCoordSetsOnGrid(tileCoordSetIter, g.RawGrid)
	withColorIter := tileSetsWithColor(onGridIter, color)
	if color == gs.ColorNone {
		return decorateUnknownNeighbors(g, withColorIter)
	}
	return withColorIter
}

// we do not iterate in any particular order since it does not matter.
// this function will only create direct paths, aka ones which would satisfy
// a Goal tile.
func (g GridSolver) dfsDirectPaths(color gs.TileColor, prev, end gs.Tile, path gs.TileCoordSet, ch chan<- gs.TileCoordSet) {

	// possible next tiles include untraversed tiles, and tiles of the same color
	possibleNext := g.RawGrid.NeighborsWith(prev.Coord, func(o gs.Tile) bool {
		return !path.Has(o.Coord) && (o.Data.Color == ColorUnknown || o.Data.Color == color)
	})

	for _, next := range possibleNext.Slice() {
		// prev's neighbors with same color, including `next`
		prevNeighborsSameColor := g.RawGrid.NeighborsWith(prev.Coord, func(o gs.Tile) bool {
			return o.Data.Color == color || path.Has(o.Coord) || o.Coord == next.Coord
		})
		// make sure that we never have an invalid path, ever
		if prevNeighborsSameColor.Len() > 2 {
			continue
		}
		if prev.Data.Type == gs.TypeGoal && prevNeighborsSameColor.Len() > 1 {
			continue
		}

		// we cannot traverse through a Goal tile that is not the end tile
		if next.Data.Type == gs.TypeGoal && next != end {
			continue
		}

		// we found a possible solution
		if next == end {

			// make sure the goal only has 1 neighbor of the same color
			if end.Data.Type == gs.TypeGoal {
				endNeighbors := g.RawGrid.NeighborsWith(end.Coord, func(o gs.Tile) bool {
					return o.Data.Color == color || path.Has(o.Coord)
				})
				if endNeighbors.Len() > 1 {
					continue
				}
			}

			var finalPath gs.TileCoordSet
			finalPath.Merge(path)
			finalPath.Add(next.Coord)
			ch <- finalPath
			continue
		}

		var nextPath gs.TileCoordSet
		nextPath.Merge(path)
		nextPath.Add(next.Coord)

		// RECURSION
		g.dfsDirectPaths(color, next, end, nextPath, ch)
	}
}

func tileCoordSetsOnGrid(coordSets <-chan gs.TileCoordSet, g gs.Grid) <-chan gs.TileSet {
	ch := make(chan gs.TileSet)
	go func() {
		for coordSet := range coordSets {
			tileSet := coordSet.ToTileSet(func(t gs.TileCoord) gs.Tile {
				return *g.TileAtCoord(t)
			})
			ch <- tileSet
		}
		close(ch)
	}()
	return ch
}
func tileSetsWithColor(tileSets <-chan gs.TileSet, color gs.TileColor) <-chan gs.TileSet {
	ch := make(chan gs.TileSet)
	go func() {
		for tileSet := range tileSets {
			var tileSetWithColor gs.TileSet
			for tile := range tileSet.Iter() {
				tileWithColor := tile
				tileWithColor.Data.Color = color
				tileSetWithColor.Add(tileWithColor)
			}
			ch <- tileSetWithColor
		}
		close(ch)
	}()
	return ch
}
func decorateUnknownNeighbors(g GridSolver, tileSets <-chan gs.TileSet) <-chan gs.TileSet {
	iter := make(chan gs.TileSet)
	go func() {
		defer close(iter)
		for tileSet := range tileSets {

			var unknownNeighbors []gs.Tile
			for tile := range tileSet.Iter() {
				neighboringUnknowns := g.RawGrid.NeighborsWith(tile.Coord, func(o gs.Tile) bool {
					return o.Data.Color == ColorUnknown && !tileSet.ToTileCoordSet().Has(o.Coord)
				})
				unknownNeighbors = append(unknownNeighbors, neighboringUnknowns.Slice()...)
			}

			for permutation := range Permutation(g.RawGrid.MaxColors-1, len(unknownNeighbors)) {
				var pathWithDecoration gs.TileSet
				pathWithDecoration.Merge(tileSet)
				for i, unknown := range unknownNeighbors {
					unknown.Data.Color = gs.TileColor(permutation[i] + 1)
					pathWithDecoration.Add(unknown)
				}
				iter <- pathWithDecoration
			}
		}
	}()
	return iter
}
