package solve

import gs "github.com/deanveloper/gridspech-go"

// FillUnknowns will look for any ColorUnknown tile touching a
// definitively ColorNone tile (ie a ColorNone path) and return
// TileSets to fill in the ColorUnknown tiles with non-ColorNone colors.
func (g GridSolver) FillUnknowns(maxColor int) <-chan gs.TileSet {
	iter := make(chan gs.TileSet)

	go func() {
		defer close(iter)

		colorNoneTiles := g.RawGrid.TilesWith(func(o gs.Tile) bool {
			return o.Data.Color == gs.ColorNone
		})

		var relevantUnknownTiles []gs.Tile
		for _, tile := range colorNoneTiles.Slice() {
			neighboringUnknowns := g.RawGrid.NeighborsWith(tile.Coord, func(o gs.Tile) bool {
				return o.Data.Color == ColorUnknown
			})
			relevantUnknownTiles = append(relevantUnknownTiles, neighboringUnknowns.Slice()...)
		}

		for permutation := range Permutation(maxColor, len(relevantUnknownTiles)) {
			var coloredUnknowns gs.TileSet
			for i, unknown := range relevantUnknownTiles {
				unknown.Data.Color = gs.TileColor(permutation[i] + 1)
				coloredUnknowns.Add(unknown)
			}
			iter <- coloredUnknowns
		}
	}()

	return iter
}
