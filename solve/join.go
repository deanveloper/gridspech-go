package solve

import (
	gs "github.com/deanveloper/gridspech-go"
)

// SolveJoins returns a channel of solutions for all of the Join tiles.
func (g GridSolver) SolveJoins() <-chan gs.TileSet {
	joinTiles := g.Grid.TilesWith(func(o gs.Tile) bool {
		return o.Data.Type == gs.TypeJoin1 || o.Data.Type == gs.TypeJoin2
	}).Slice()

	if len(joinTiles) == 0 {
		ch := make(chan gs.TileSet, 1)
		ch <- gs.NewTileSet()
		close(ch)
		return ch
	}

	tilesToSolutions := make([]<-chan gs.TileSet, len(joinTiles))
	for i, tile := range joinTiles {
		tilesToSolutions[i] = g.SolveJoin(tile)
	}

	// now merge them all together
	for i := 1; i < len(joinTiles); i++ {
		mergedIter := mergeSolutionsIters(tilesToSolutions[i-1], tilesToSolutions[i])
		tilesToSolutions[i] = mergedIter
	}

	return tilesToSolutions[len(tilesToSolutions)-1]
}

// SolveJoin returns a channel of solutions for an individual join tile.
func (g GridSolver) SolveJoin(join gs.Tile) <-chan gs.TileSet {
	joinIter := make(chan gs.TileSet)

	go func() {
		defer close(joinIter)

		var joinNum int
		switch join.Data.Type {
		case gs.TypeJoin1:
			joinNum = 1
		case gs.TypeJoin2:
			joinNum = 2
		default:
			panic("not a join tile")
		}

		for c := 0; c < g.Grid.MaxColors; c++ {
			color := gs.TileColor(c)
			shapeCh, pruneCh := g.ShapesIter(join.Coord, color)
			for shape := range shapeCh {
				prune := shouldPruneJoin(g, shape, color, joinNum)
				pruneCh <- prune

				specialTiles := numSpecialTiles(g, shape, joinNum)
				if !prune && specialTiles == joinNum+1 {
					for decorated := range decorateSetBorder(g, color, shape) {
						joinIter <- decorated
					}
				}
			}
		}
	}()

	return joinIter
}

// trim if:
// - too many special tiles in the shape
// - if the shape contains a goal tile, it must be a path
// - if joinNum is 1, it cannot contain a goal tile
func shouldPruneJoin(g GridSolver, shape gs.TileSet, color gs.TileColor, joinNum int) bool {
	shapeAsSlice := shape.Slice()
	shapeCoords := shape.ToTileCoordSet()

	var containsGoalTile bool
	var containsTrineighborTile bool
	var specialTiles int
	for _, tile := range shapeAsSlice {

		if tile.Data.Type != gs.TypeBlank {
			specialTiles++
			if specialTiles > joinNum+1 {
				return true
			}
		}

		sameColorNeighborsInShape := g.Grid.NeighborSetWith(tile.Coord, func(o gs.Tile) bool {
			return shapeCoords.Has(o.Coord)
		})
		if sameColorNeighborsInShape.Len() > 2 {
			containsTrineighborTile = true
		}
		if tile.Data.Type == gs.TypeGoal {
			if joinNum == 1 {
				return true
			}
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

func numSpecialTiles(g GridSolver, shape gs.TileSet, joinNum int) int {
	var numSpecialTiles int
	for _, tile := range shape.Slice() {
		if tile.Data.Type != gs.TypeBlank {
			numSpecialTiles++
		}
	}
	return numSpecialTiles
}
