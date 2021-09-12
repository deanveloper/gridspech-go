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

	tilesToSolutions := make([]<-chan gs.TileSet, len(crownTiles))
	for i, tile := range crownTiles {
		tilesToSolutions[i] = g.solveCrown(tile.Coord)
	}

	// now merge them all together
	for i := 1; i < len(crownTiles); i++ {
		mergedIter := mergeSolutionsIters(tilesToSolutions[i-1], tilesToSolutions[i])
		tilesToSolutions[i] = mergedIter
	}

	return tilesToSolutions[len(tilesToSolutions)-1]
}

func (g GridSolver) solveCrown(crown gs.TileCoord) <-chan gs.TileSet {

	crownIter := make(chan gs.TileSet)

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
// - this list used to be longer lmao
func shouldPruneCrown(g GridSolver, crown gs.TileCoord, shape gs.TileSet, color gs.TileColor) bool {

	for _, tile := range shape.Slice() {
		if tile.Data.Type == gs.TypeCrown && tile.Data.Color == color && tile.Coord != crown {
			return true
		}
	}

	return false
}

func filterValidSoFar(
	g GridSolver,
	previousTiles []gs.Tile,
	current gs.Tile,
	sols <-chan gs.TileSet,
) <-chan gs.TileSet {
	filtered := make(chan gs.TileSet, 200)

	go func() {
		defer close(filtered)
		for solution := range sols {
			newBase := g.Grid.Clone()
			newBase.ApplyTileSet(solution)

			allValid := true
			for _, tile := range previousTiles {
				if !newBase.ValidTile(tile.Coord) {
					allValid = false
					break
				}
			}
			if allValid {
				filtered <- solution
			}
		}
	}()

	return filtered
}
