package solve

import (
	gs "github.com/deanveloper/gridspech-go"
)

// PathsIter returns an channel of direct paths from start to end.
// These paths will:
//   1. never contain a goal tile that isn't start or end.
//   2. never make a path that would cause start or end to become invalid Goal tiles.
//   3. have the same Color as start.
func (g GridSolver) PathsIter(start, end gs.TileCoord, color gs.TileColor) <-chan gs.TileSet {
	tileCoordSetIter := make(chan gs.TileCoordSet)
	go func() {
		defer close(tileCoordSetIter)
		startTile, endTile := *g.Grid.TileAtCoord(start), *g.Grid.TileAtCoord(end)
		if !g.UnknownTiles.Has(start) && color != startTile.Data.Color {
			return
		}
		if !g.UnknownTiles.Has(end) && color != endTile.Data.Color {
			return
		}
		g.dfsDirectPaths(color, startTile, endTile, gs.NewTileCoordSet(start), tileCoordSetIter)
	}()

	onGridIter := decorateWithGridInfo(tileCoordSetIter, g.Grid)
	withColorIter := decorateWithColor(onGridIter, color)
	if color == gs.ColorNone {
		return decorateUnknownNeighbors(g, withColorIter)
	}
	return withColorIter
}

// we do not iterate in any particular order since it does not matter.
// this function will only create direct paths, aka ones which would satisfy
// a Goal tile.
func (g GridSolver) dfsDirectPaths(color gs.TileColor, prev, end gs.Tile, path gs.TileCoordSet, ch chan<- gs.TileCoordSet) {

	// possible next tiles include unknown tiles, and tiles of the same color
	possibleNext := g.Grid.NeighborsWith(prev.Coord, func(o gs.Tile) bool {
		return !path.Has(o.Coord) && (g.UnknownTiles.Has(o.Coord) || o.Data.Color == color)
	})

	for _, next := range possibleNext.Slice() {
		// prev's neighbors we _know_ are same color (including those that are part of the path)
		prevNeighborsSameColor := g.Grid.NeighborsWith(prev.Coord, func(o gs.Tile) bool {
			knownSameColor := (o.Data.Color == color && !g.UnknownTiles.Has(o.Coord))
			partOfPath := path.Has(o.Coord) || o.Coord == next.Coord
			return knownSameColor || partOfPath
		})
		// make sure that we never have an invalid path
		if prevNeighborsSameColor.Len() > 2 {
			continue
		}
		if prev.Data.Type == gs.TypeGoal && prevNeighborsSameColor.Len() > 1 {
			continue
		}

		if next.Data.Type == gs.TypeGoal && next.Coord != end.Coord {
			continue
		}

		// we found a possible solution
		if next == end {

			// make sure the goal only has 1 neighbor we know is the same color
			if end.Data.Type == gs.TypeGoal {
				endNeighbors := g.Grid.NeighborsWith(end.Coord, func(o gs.Tile) bool {
					return (o.Data.Color == color && !g.UnknownTiles.Has(o.Coord)) || path.Has(o.Coord)
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

func decorateWithGridInfo(coordSets <-chan gs.TileCoordSet, g gs.Grid) <-chan gs.TileSet {
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
func decorateWithColor(tileSets <-chan gs.TileSet, color gs.TileColor) <-chan gs.TileSet {
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
				neighboringUnknowns := g.Grid.NeighborsWith(tile.Coord, func(o gs.Tile) bool {
					return g.UnknownTiles.Has(o.Coord) && !tileSet.ToTileCoordSet().Has(o.Coord)
				})
				unknownNeighbors = append(unknownNeighbors, neighboringUnknowns.Slice()...)
			}

			for permutation := range Permutation(g.Grid.MaxColors-1, len(unknownNeighbors)) {
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
