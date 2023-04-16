package solve

import (
	gs "github.com/deanveloper/gridspech-go"
)

// SolveCrowns will return a channel of solutions for all the crown tiles in g.
func (g GridSolver) SolveCrowns() <-chan gs.TileSet {

	// get all crown tiles
	crownTiles := g.Grid.TilesWith(func(o gs.Tile) bool {
		return o.Data.Type == gs.TypeCrown
	}).Slice()

	if len(crownTiles) == 0 {
		ch := make(chan gs.TileSet, 1)
		ch <- gs.NewTileSet()
		close(ch)
		return ch
	}

	tilesToSolutions := make([]<-chan gs.TileSet, len(crownTiles))
	for i, tile := range crownTiles {
		tilesToSolutions[i] = g.SolveCrown(tile.Coord)
	}

	// now merge them all together
	for i := 1; i < len(crownTiles); i++ {
		mergedIter := MergeSolutionsIters(tilesToSolutions[i-1], tilesToSolutions[i])
		tilesToSolutions[i] = mergedIter
	}

	return tilesToSolutions[len(tilesToSolutions)-1]
}

// SolveCrown returns a channel of solutions for a crown at the given coordinate.
func (g GridSolver) SolveCrown(crown gs.TileCoord) <-chan gs.TileSet {

	crownIter := make(chan gs.TileSet, 50)

	go func() {
		defer close(crownIter)

		for c := 0; c < g.Grid.MaxColors; c++ {

			shapesCh, pruneCh := g.ShapesIter(crown, gs.TileColor(c))

			for shape := range shapesCh {
				prune := shouldPruneCrown(g, crown, shape, gs.TileColor(c))
				pruneCh <- prune
				if !prune {
					for decorated := range decorateSetBorder(g, gs.TileColor(c), shape) {
						crownIter <- decorated
					}
				}
			}
		}
	}()

	return crownIter
}

// prune if:
// - this shape contains a separate crown of the same color
// - if there is a goal tile, the shape must be a path
func shouldPruneCrown(g GridSolver, crown gs.TileCoord, shape gs.TileSet, color gs.TileColor) bool {
	shapeAsSlice := shape.Slice()
	shapeCoords := shape.ToTileCoordSet()

	var containsGoalTile bool
	var containsTrineighborTile bool
	for _, tile := range shapeAsSlice {
		if tile.Data.Type == gs.TypeCrown && tile.Coord != crown {
			return true
		}

		sameColorNeighborsInShape := g.Grid.NeighborSetWith(tile.Coord, func(o gs.Tile) bool {
			return shapeCoords.Has(o.Coord)
		})
		if sameColorNeighborsInShape.Len() > 2 {
			containsTrineighborTile = true
		}
		if tile.Data.Type == gs.TypeGoal {
			if sameColorNeighborsInShape.Len() > 1 {
				return true
			}
			containsGoalTile = true
		}

		if containsGoalTile && containsTrineighborTile {
			return true
		}
	}

	return false
}
