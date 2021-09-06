package solve

import (
	gs "github.com/deanveloper/gridspech-go"
)

// SolveCrowns will return a channel of solutions for all the crown tiles in g.
//
// For performance reasons, it is expected that you run this after all other tiles have been solved.
func (g GridSolver) SolveCrowns() <-chan gs.TileSet {
	return nil
}

func (g GridSolver) solveCrown(crown gs.TileCoord, color gs.TileColor) <-chan gs.TileSet {

	crownIter := make(chan gs.TileSet)

	go func() {

		shapesCh, pruneCh := g.ShapesIter(crown, gs.TileColor(color))

		for shape := range shapesCh {
			prune := shouldPrune(g, crown, shape, color)
			pruneCh <- prune
			if !prune {
				crownIter <- shape
			}
		}

		close(crownIter)
	}()

	return crownIter
}

// prune if:
// - this shape contains a separate crown of the same color
// - has a goal tile, and:
//   - any tile has >2 neighbors
//   - any goal tile has >1 neighbor
func shouldPrune(g GridSolver, crown gs.TileCoord, shape gs.TileSet, color gs.TileColor) bool {
	shapeAsSlice := shape.Slice()
	var coordsWithColor gs.TileCoordSet
	for _, tile := range shapeAsSlice {
		if tile.Data.Color == color {
			coordsWithColor.Add(tile.Coord)
		}
	}

	var containsGoalTile bool
	var containsTrineighborTile bool
	for _, tile := range shapeAsSlice {
		if tile.Data.Type == gs.TypeCrown && tile.Data.Color == color && tile.Coord != crown {
			return true
		}

		sameColorNeighborsInShape := g.Grid.NeighborsWith(tile.Coord, func(o gs.Tile) bool {
			return coordsWithColor.Has(o.Coord)
		})
		if sameColorNeighborsInShape.Len() > 2 {
			containsTrineighborTile = true
		}
		if tile.Data.Type == gs.TypeGoal && tile.Data.Color == color {
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
