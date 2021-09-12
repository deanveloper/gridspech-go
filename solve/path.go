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
	pathIter := make(chan gs.TileSet)
	go func() {
		defer close(pathIter)
		startTile, endTile := *g.Grid.TileAtCoord(start), *g.Grid.TileAtCoord(end)
		if !g.UnknownTiles.Has(start) && color != startTile.Data.Color {
			return
		}
		if !g.UnknownTiles.Has(end) && color != endTile.Data.Color {
			return
		}
		g.dfsDirectPaths(color, startTile, endTile, gs.NewTileCoordSet(start), pathIter)
	}()

	withBorderIter := decorateSetIterBorders(g, color, pathIter)
	return withBorderIter
}

// we do not iterate in any particular order since it does not matter.
// this function will only create direct paths, aka ones which would satisfy
// a Goal tile.
func (g GridSolver) dfsDirectPaths(color gs.TileColor, prev, end gs.Tile, path gs.TileCoordSet, ch chan<- gs.TileSet) {

	// possible next tiles include unknown tiles, and tiles of the target color
	possibleNext := g.Grid.NeighborSetWith(prev.Coord, func(o gs.Tile) bool {
		return !path.Has(o.Coord) && (g.UnknownTiles.Has(o.Coord) || o.Data.Color == color)
	})

	for _, next := range possibleNext.Slice() {
		// prev's neighbors we _know_ are same color (including those that are part of the path)
		prevNeighborsSameColor := g.Grid.NeighborSetWith(prev.Coord, func(o gs.Tile) bool {
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
				endNeighbors := g.Grid.NeighborSetWith(end.Coord, func(o gs.Tile) bool {
					return (o.Data.Color == color && !g.UnknownTiles.Has(o.Coord)) || path.Has(o.Coord)
				})
				if endNeighbors.Len() > 1 {
					continue
				}
			}

			finalPath := path.ToTileSet(func(t gs.TileCoord) gs.Tile {
				tileCopy := *g.Grid.TileAtCoord(t)
				tileCopy.Data.Color = color
				return tileCopy
			})
			next.Data.Color = color
			finalPath.Add(next)
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
