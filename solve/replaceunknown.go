package solve

import gs "github.com/deanveloper/gridspech-go"

// FillUnknowns will look for any ColorUnknown tile touching a
// definitively ColorNone tile (ie a ColorNone path) and return
// TileSets to fill in the ColorUnknown tiles with non-ColorNone colors.
func (g GridSolver) FillUnknowns(maxColor int) func(yield func(ts gs.TileSet)) {
	colorNoneTiles := g.RawGrid.TilesWith(func(o gs.Tile) bool {
		return o.Color == gs.ColorNone
	})

	var relevantUnknownTiles []gs.Tile
	for _, tile := range colorNoneTiles.Slice() {
		neighboringUnknowns := g.RawGrid.NeighborsWith(tile, func(o gs.Tile) bool {
			return o.Color == ColorUnknown
		})
		relevantUnknownTiles = append(relevantUnknownTiles, neighboringUnknowns.Slice()...)
	}

	return func(yield func(ts gs.TileSet)) {
		permIter := Permutation(maxColor, len(relevantUnknownTiles))
		permIter(func(permutation []int) {
			var coloredUnknowns gs.TileSet
			for i, unknown := range relevantUnknownTiles {
				unknown.Color = gs.TileColor(permutation[i] + 1)
				coloredUnknowns.Add(unknown)
			}
			yield(coloredUnknowns)
		})
	}
}
