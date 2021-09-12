package solve

import (
	gs "github.com/deanveloper/gridspech-go"
)

// SolveAllTiles returns a channel which will return a TileSet of all tiles in g.
func (g GridSolver) SolveAllTiles() <-chan gs.TileSet {
	solutionIter := make(chan gs.TileSet)

	go func() {
		defer close(solutionIter)

		for goalsAndDots := range mergeSolutionsIters(g.SolveGoals(), g.SolveDots()) {
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
