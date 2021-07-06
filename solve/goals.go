package solve

import (
	"fmt"
	"log"
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

// allows us to use a slice of goal pairings in a map
type goalPairingKey string

func makeGoalPairingKey(pairings [][2]gs.Tile) goalPairingKey {
	return goalPairingKey(fmt.Sprint(pairings))
}

// Goals will return a channels of solutions for all the goal tiles in g
func Goals(g Grid, maxColors gs.TileColor) <-chan Grid {

	var gridChanWg sync.WaitGroup
	gridChan := make(chan Grid, 50)

	// get all goal tiles
	var goalTiles []gs.Tile
	for _, col := range g.Tiles {
		for _, tile := range col {
			if tile.Type == gridspech.TypeGoal {
				goalTiles = append(goalTiles, tile)
			}
		}
	}

	goalPairingsSet := allTilePairingSets(goalTiles)

	// map of pairs to pairings which contain that pair
	pairToPairingsSet := make(map[[2]gs.Tile]([][][2]gs.Tile))
	for _, pairings := range goalPairingsSet {
		for _, pairing := range pairings {
			pairToPairingsSet[pairing] = append(pairToPairingsSet[pairing], pairings)
		}
	}

	// map of pair to solutions (so far) for that pair
	var pairToSolutionsLock sync.RWMutex
	pairToSolutions := make(map[[2]gs.Tile][]goalSolution)
	for k := range pairToPairingsSet {
		for color := gs.ColorNone; color < maxColors; color++ {
			k, color := k, color
			pair := [2]gs.Tile{k[0], k[1]}
			gridChanWg.Add(1)
			go func() {
				for solution := range g.SolvePath(k[0], k[1], color) {
					newGoalSolution := goalSolution{
						start:    k[0],
						end:      k[1],
						color:    color,
						solution: solution,
					}

					pairToSolutionsLock.RLock()
					grids := onNewSolutionFound(g, newGoalSolution, pairToPairingsSet[pair], pairToSolutions)
					pairToSolutionsLock.RUnlock()

					pairToSolutionsLock.Lock()
					pairToSolutions[pair] = append(pairToSolutions[pair], newGoalSolution)
					pairToSolutionsLock.Unlock()

					for _, grid := range grids {
						invalidGoals := grid.TilesWith(func(t gs.Tile) bool {
							return t.Type == gs.TypeGoal && !grid.ValidTile(t)
						})
						if invalidGoals.Len() > 0 {
							continue
						}
						gridChan <- grid
					}
				}
				gridChanWg.Done()
			}()
		}
	}

	go func() {
		gridChanWg.Wait()
		close(gridChan)
	}()

	return gridChan
}

func onNewSolutionFound(
	baseGrid Grid,
	newSolution goalSolution,
	pairingsToUpdate [][][2]gs.Tile,
	currentSolutions map[[2]gs.Tile][]goalSolution,
) []Grid {
	var grids []Grid
pairingsToUpdateLoop:
	for _, pairings := range pairingsToUpdate {
		solutions := make(map[[2]gs.Tile][]goalSolution, len(pairings)-1)
		for _, pair := range pairings {
			newPair := [2]gs.Tile{newSolution.start, newSolution.end}
			if pair != newPair {
				// if this pair does not have any solutions yet,
				// we do not care about this pairing
				if len(solutions) == 0 {
					continue pairingsToUpdateLoop
				}
				solutions[pair] = currentSolutions[pair]
			} else {
				solutions[pair] = []goalSolution{newSolution}
			}
		}

		forEachSolutionSet(solutions, func(gs []goalSolution) {
			grids = append(grids, combineSolutions(baseGrid, gs))
		})
	}

	return grids
}

func forEachSolutionSet(solutionSet map[[2]gs.Tile][]goalSolution, forEach func([]goalSolution)) {
	for pair, solutions := range solutionSet {

		remainingSolutions := make(map[[2]gs.Tile][]goalSolution, len(solutionSet)-1)
		for pair2, solutions2 := range solutionSet {
			if pair2 != pair {
				remainingSolutions[pair] = solutions2
			}
		}

		if len(remainingSolutions) == 0 {
			for _, solution := range solutions {
				forEach([]goalSolution{solution})
			}
		} else {
			for _, solution := range solutions {
				forEachSolutionSet(remainingSolutions, func(solSet []goalSolution) {
					newGoalSolution := make([]goalSolution, len(solSet)+1)
					newGoalSolution[0] = solution
					copy(newGoalSolution[1:], solSet)
					if anyIntersections(newGoalSolution) {
						return
					}

					forEach(newGoalSolution)
				})
			}
		}
	}
}

func anyIntersections(solSet []goalSolution) bool {
	var allTiles gs.TileSet

	for _, solution := range solSet {
		for _, tile := range solution.solution.Slice() {
			if allTiles.Has(tile) {
				return true
			}
			allTiles.Add(tile)
		}
	}

	return false
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

				// and just in case... validate all goals in the grid!
				invalidGoals := grid.TilesWith(func(t gs.Tile) bool {
					return t.Type == gs.TypeGoal && !grid.ValidTile(t)
				})
				if invalidGoals.Len() > 0 {
					log.Printf("invalid goals found %v", invalidGoals)
					continue
				}

				ch <- grid
			}

			// add this solution to previously found solutions
			goalsToSolutions[solution.start] = append(goalsToSolutions[solution.start], solution)
		}
		close(ch)
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

		// skip if relevant pairing set
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
			for _, oldSolution := range solutionsForGoals[pairing[0]] {
				if oldSolution.end == pairing[1] {
					goalsForPairing = append(goalsForPairing, oldSolution)
				}
			}
			pairingToSolutions[pairing] = goalsForPairing
		}
		solPairing := [2]gs.Tile{solution.start, solution.end}
		pairingToSolutions[solPairing] = append(pairingToSolutions[solPairing], solution)

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
		if len(remainingSolutions) == 0 {
			result = append(result, []goalSolution{newSolution})
		}
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

	pairingSets := AllPairingSets(len(tiles))
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

	for _, sol := range sols {
		for _, tile := range sol.solution.Slice() {
			grid.Tiles[tile.X][tile.Y].Color = sol.color
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
				start:    solsCh.start,
				end:      solsCh.end,
				color:    solsCh.color,
				solution: sol,
			}
		}
		wg.Done()
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
