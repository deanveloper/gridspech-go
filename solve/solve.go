package solve

import (
	"fmt"

	gs "github.com/deanveloper/gridspech-go"
)

// SolveAllTiles returns a channel which will return a TileSet of all tiles in g.
func (g GridSolver) SolveAllTiles() <-chan gs.TileSet {
	solutionIter := make(chan gs.TileSet)

	go func() {
		defer close(solutionIter)

		for goalsAndDots := range MergeSolutionsIters(g.SolveGoals(), g.SolveDots()) {
			newGrid := g.Clone()
			newGrid.Grid.ApplyTileSet(goalsAndDots)
			newGrid.UnknownTiles.RemoveAll(goalsAndDots.ToTileCoordSet())

			for joinsSolution := range newGrid.SolveJoins() {
				joinsSolved := newGrid.Clone()
				joinsSolved.Grid.ApplyTileSet(joinsSolution)
				joinsSolved.UnknownTiles.RemoveAll(joinsSolution.ToTileCoordSet())

				for crownsSolution := range joinsSolved.SolveCrowns() {
					crownsSolved := joinsSolved.Clone()
					crownsSolved.Grid.ApplyTileSet(crownsSolution)
					crownsSolved.UnknownTiles.RemoveAll(crownsSolution.ToTileCoordSet())

					if crownsSolved.Grid.Valid() {
						var merged gs.TileSet
						merged.Merge(goalsAndDots)
						merged.Merge(joinsSolution)
						merged.Merge(crownsSolution)
						solutionIter <- merged
					}
				}
			}
		}
	}()

	return solutionIter
}

// SolveTiles returns a channel of possible solutions for the given tiles.
func (g GridSolver) SolveTiles(tiles ...gs.TileCoord) <-chan gs.TileSet {

	if len(tiles) == 0 {
		ch := make(chan gs.TileSet, 1)
		ch <- gs.NewTileSet()
		close(ch)
		return ch
	}

	tilesToSolutions := make([]<-chan gs.TileSet, len(tiles))
	for i, tile := range tiles {
		tilesToSolutions[i] = g.solveTile(*g.Grid.TileAtCoord(tile))
	}

	// now merge them all together
	for i := 1; i < len(tiles); i++ {
		mergedIter := MergeSolutionsIters(tilesToSolutions[i-1], tilesToSolutions[i])
		uniqueIter := filterUnique(mergedIter)
		tilesToSolutions[i] = uniqueIter
	}

	return tilesToSolutions[len(tilesToSolutions)-1]
}

func (g GridSolver) solveTile(t gs.Tile) <-chan gs.TileSet {
	switch t.Data.Type {
	case gs.TypeHole, gs.TypeBlank:
		ch := make(chan gs.TileSet, 1)
		ch <- gs.NewTileSet()
		close(ch)
		return ch
	case gs.TypeGoal:
		return filterHasTile(g.SolveGoals(), t.Coord)
	case gs.TypeCrown:
		return g.SolveCrown(t.Coord)
	case gs.TypeDot1, gs.TypeDot2, gs.TypeDot3:
		return g.SolveDot(t)
	case gs.TypeJoin1, gs.TypeJoin2:
		return g.SolveJoin(t)
	default:
		panic(fmt.Sprintf("invalid type %v", t.Data.Type))
	}
}
