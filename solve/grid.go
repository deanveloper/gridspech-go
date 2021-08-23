package solve

import gs "github.com/deanveloper/gridspech-go"

// GridSolver represents a gs.Grid, but with a special "unknown" tile color.
type GridSolver struct {
	Grid         gs.Grid
	UnknownTiles gs.TileCoordSet
}

// NewGridSolver creates a GridSolver
func NewGridSolver(solving gs.Grid) GridSolver {
	newSolving := solving.Clone()
	nonSticky := newSolving.TilesWith(func(o gs.Tile) bool {
		return !o.Data.Sticky
	})
	var unknownTiles gs.TileCoordSet
	for tile := range nonSticky.Iter() {
		unknownTiles.Add(tile.Coord)
	}
	return GridSolver{Grid: newSolving, UnknownTiles: unknownTiles}
}

// Clone returns clone of g.
func (g GridSolver) Clone() GridSolver {
	newUnknownTiles := gs.NewTileCoordSet()
	newUnknownTiles.Merge(g.UnknownTiles)
	return GridSolver{Grid: g.Grid.Clone(), UnknownTiles: newUnknownTiles}
}
