package solve

import (
	"sync"

	"github.com/deanveloper/gridspech-go"
	gs "github.com/deanveloper/gridspech-go"
)

type goalSolutionsChan struct {
	start, end gs.Tile
	color      gs.TileColor
	sols       <-chan gs.TileSet
}

type goalSolution struct {
	start, end gs.Tile
	color      gs.TileColor
	solution   gs.TileSet
}

func (g goalSolution) eq(o goalSolution) bool {
	return g.start == o.start && g.end == o.end && g.color == o.color && g.solution.Eq(o.solution)
}

// Goals will return a channels of solutions for all the goal tiles in g
func Goals(g Grid, maxColors gs.TileColor) <-chan Grid {

	// get all goal tiles
	var goalTiles []gs.Tile
	for _, col := range g.Tiles {
		for _, tile := range col {
			if tile.Type == gridspech.TypeGoal {
				goalTiles = append(goalTiles, tile)
			}
		}
	}

	// find all solutions for each goal tile going to every other goal tile
	var solutionsChans []goalSolutionsChan
	for i1 := 0; i1 < len(goalTiles)-1; i1++ {
		for i2 := i1 + 1; i2 < len(goalTiles); i2++ {
			goal1, goal2 := goalTiles[i1], goalTiles[i2]

			for color := gs.TileColor(0); color < maxColors; color++ {
				paths := g.SolvePath(goal1, goal2, color)
				solutionsChans = append(solutionsChans, goalSolutionsChan{
					start: goal1,
					end:   goal2,
					color: color,
					sols:  paths,
				})
			}
		}
	}

	// aggregate solutions into a single channel
	allSolutions := aggregateGoalSolutionChans(solutionsChans)

	// make a map of tiles to their solutions, to be populated as time goes on
	return assembleSolutions(g, allSolutions, goalTiles)
}

// assembles solutions together
func assembleSolutions(baseGrid Grid, allSolutions <-chan goalSolution, goalTiles []gs.Tile) <-chan Grid {

	ch := make(chan Grid)
	go func() {
		goalsToSolutions := make(map[gs.Tile][]goalSolution)
		pairingsSet := allTilePairingSets(goalTiles)

		for solution := range allSolutions {

			// first, try out the solution with all prev solutions found
			grids := makeFullSolutions(baseGrid, solution, goalsToSolutions, pairingsSet)
			for _, grid := range grids {
				ch <- grid
			}

			// add this solution to previously found solutions
			goalsToSolutions[solution.start] = append(goalsToSolutions[solution.start], solution)
		}
	}()

	return ch
}

func makeFullSolutions(
	baseGrid Grid,
	solution goalSolution,
	solutionsForGoals map[gs.Tile][]goalSolution,
	pairingsSet [][][2]gs.Tile,
) []Grid {

	var grids []Grid
	for _, pairings := range pairingsSet {
		var relevantPairingSet bool
		for _, pairing := range pairings {
			if pairing[0] == solution.start && pairing[1] == solution.end {
				relevantPairingSet = true
				break
			}
		}
		if !relevantPairingSet {
			continue
		}

		pairingToSolutions := make(map[[2]gs.Tile][]goalSolution)
		for _, pairing := range pairings {
			var goalsForPairing []goalSolution
			for _, solution := range solutionsForGoals[pairing[0]] {
				if solution.end == pairing[1] {
					goalsForPairing = append(goalsForPairing, solution)
				}
			}
			pairingToSolutions[pairing] = goalsForPairing
		}

		eachSolutionSet := eachSolutionForPairings(pairings, pairingToSolutions)
		for _, solutionSet := range eachSolutionSet {
			grids = append(grids, combineSolutions(baseGrid, solutionSet))
		}
	}
	return grids
}

func eachSolutionForPairings(pairings [][2]gs.Tile, pairingToSolutions map[[2]gs.Tile][]goalSolution) [][]goalSolution {
	if len(pairings) == 0 {
		return nil
	}

	pairing := pairings[0]
	solutions := pairingToSolutions[pairing]

	remainingSolutions := eachSolutionForPairings(pairings[1:], pairingToSolutions)

	var result [][]goalSolution
	for _, newSolution := range solutions {
	oldSolutionsLoop:
		for _, oldSolutions := range remainingSolutions {
			newSolutions := make([]goalSolution, len(oldSolutions)+1)
			newSolutions[0] = newSolution

			var clone gs.TileSet
			clone.Merge(newSolution.solution)
			for i := range oldSolutions {
				oldSolution := oldSolutions[i]

				// while adding solutions, check to make sure that the
				// new solution does not intersect with any current ones.
				clone.RemoveAll(oldSolution.solution)
				if !clone.Eq(newSolution.solution) {
					continue oldSolutionsLoop
				}
				newSolutions[i+1] = oldSolution
			}
			result = append(result, newSolutions)
		}
	}

	return result
}

func allTilePairingSets(tiles []gs.Tile) [][][2]gs.Tile {
	ints := make([]int, len(tiles))
	for i := range ints {
		ints[i] = i
	}

	pairingSets := AllPairingSets(ints)
	tilePairingSets := make([][][2]gs.Tile, len(pairingSets))
	for i, pairings := range pairingSets {
		tilePairings := make([][2]gs.Tile, len(pairings))
		for p, pairing := range pairings {
			tilePairings[p] = [2]gs.Tile{tiles[pairing[0]], tiles[pairing[1]]}
		}
		tilePairingSets[i] = tilePairings
	}
	return tilePairingSets
}

// combines solutions together and returns the grid of the combined solutions.
// Second value is false if the solutions could not be combined.
func combineSolutions(baseGrid Grid, sols []goalSolution) Grid {
	grid := baseGrid.Clone()

	for i1 := 0; i1 < len(sols)-1; i1++ {
		for _, tile := range sols[i1].solution.Slice() {
			grid.Tiles[tile.X][tile.Y].Color = sols[i1].color
		}
	}

	return Grid{Grid: grid}
}

func aggregateGoalSolutionChans(solutions []goalSolutionsChan) <-chan goalSolution {
	ch := make(chan goalSolution)

	var wg sync.WaitGroup
	wg.Add(len(solutions))

	mergeIn := func(solsCh goalSolutionsChan) {
		for sol := range solsCh.sols {
			ch <- goalSolution{
				color:    solsCh.color,
				solution: sol,
			}
			wg.Done()
		}
	}

	for _, solutionCh := range solutions {
		go mergeIn(solutionCh)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}
