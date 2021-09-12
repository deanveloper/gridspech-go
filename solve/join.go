package solve

import gs "github.com/deanveloper/gridspech-go"

// SolveJoins returns a channel of solutions for all of the Join tiles.
func (g GridSolver) SolveJoins() <-chan gs.TileSet {
	joinTiles := g.Grid.TilesWith(func(o gs.Tile) bool {
		return o.Data.Type == gs.TypeJoin1 || o.Data.Type == gs.TypeJoin2
	}).Slice()

	tilesToSolutions := make([]<-chan gs.TileSet, len(joinTiles))
	for i, tile := range joinTiles {
		tilesToSolutions[i] = g.SolveCrown(tile.Coord)
	}

	// now merge them all together
	for i := 1; i < len(joinTiles); i++ {
		mergedIter := mergeSolutionsIters(tilesToSolutions[i-1], tilesToSolutions[i])
		tilesToSolutions[i] = mergedIter
	}

	return nil
}

// SolveJoin returns a channel of solutions for an individual join tile.
func (g GridSolver) SolveJoin(join gs.Tile) <-chan gs.TileSet {
	joinIter := make(chan gs.TileSet)

	go func() {
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
			shapeCh, pruneCh := g.ShapesIter(join.Coord, gs.TileColor(c))
			for shape := range shapeCh {
				specialTiles := numSpecialTiles(g, shape, joinNum)
				pruneCh <- (specialTiles > joinNum+1)

				if specialTiles == joinNum+1 {
					for decorated := range decorateSetBorder(g, gs.TileColor(c), shape) {
						joinIter <- decorated
					}
				}
			}
		}
	}()

	return joinIter
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
