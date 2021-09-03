package solve

import gs "github.com/deanveloper/gridspech-go"

// Kings will return a channel of solutions for all the king tiles in g.
func Kings(g GridSolver) <-chan gs.TileSet {
	return nil
}

// crowns are very tough to solve, so we basically just have to "abort" if any tiles which are not crowns become invalid.
func (g GridSolver) solveKing() {
}

// SolveShapes returns an iterator of all shapes which contain `start`, and be made out `color`.
func (g GridSolver) SolveShapes(start gs.TileCoord, color gs.TileColor) <-chan gs.TileSet {
	ch := make(chan gs.TileSet)

	go func() {
		g.solveShapesRecur(start, color, gs.NewTileCoordSet(), ch)
		close(ch)
	}()

	return ch
}

func (g GridSolver) solveShapesRecur(current gs.TileCoord, color gs.TileColor, shape gs.TileCoordSet, ch chan<- gs.TileSet) {
	shape.Add(current)

	ch <- shape.ToTileSet(func(t gs.TileCoord) gs.Tile {
		tileCopy := *g.Grid.TileAtCoord(t)
		tileCopy.Data.Color = color
		return tileCopy
	})

	nextTiles := g.Grid.NeighborsWith(current, func(o gs.Tile) bool {
		return !shape.Has(o.Coord) && (g.UnknownTiles.Has(o.Coord) || o.Data.Color == color)
	})
	for next := range nextTiles.Iter() {
		g.solveShapesRecur(next.Coord, color, shape, ch)
	}

	shape.Remove(current)
}
