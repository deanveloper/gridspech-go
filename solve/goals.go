package solve

import (
	"fmt"
	"sync"

	"github.com/deanveloper/gridspech-go"
	gs "github.com/deanveloper/gridspech-go"
)

type goalSolution struct {
	start, end gs.Tile
	color      gs.TileColor
	path       gs.TileSet
}

// allows us to use a slice of goal pairings in a map
type goalPairingKey string

func makeGoalPairingKey(pairings [][2]gs.Tile) goalPairingKey {
	return goalPairingKey(fmt.Sprint(pairings))
}

// Goals will return a channel of solutions for all the goal tiles in g
func Goals(g GridSolver, maxColors int) <-chan gs.TileSet {

	var finalSolutionsWg sync.WaitGroup
	finalSolutionsCh := make(chan gs.TileSet, 5)

	// get all goal tiles
	var goalTiles []gs.Tile
	for _, col := range g.RawGrid.Tiles {
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
	var pairToSolutionsLock sync.Mutex
	pairToSolutions := make(map[[2]gs.Tile][]goalSolution)
	for k := range pairToPairingsSet {
		for color := 0; color < maxColors; color++ {
			k, color := k, color
			pair := [2]gs.Tile{k[0], k[1]}
			finalSolutionsWg.Add(1)
			go func() {
				for solution := range g.SolvePath(k[0], k[1], gs.TileColor(color)) {

					var solutionWithColor gs.TileSet
					for _, tile := range solution.Slice() {
						tileWithColor := tile
						tileWithColor.Color = gs.TileColor(color)
						solutionWithColor.Add(tileWithColor)
					}

					newGoalSolution := goalSolution{
						start: k[0],
						end:   k[1],
						color: gs.TileColor(color),
						path:  solutionWithColor,
					}

					pairToSolutionsLock.Lock()
					newSolutionSets := findNewSolutionSets(g, newGoalSolution, pairToPairingsSet[pair], pairToSolutions)

					pairToSolutions[pair] = append(pairToSolutions[pair], newGoalSolution)
					pairToSolutionsLock.Unlock()

					for _, solutionSet := range newSolutionSets {
						// make sure that we don't make any invalid goals
						grid := g.Clone()
						grid.RawGrid.ApplyTileSet(solutionSet)
						invalidGoals := grid.RawGrid.TilesWith(func(t gs.Tile) bool {
							return t.Type == gs.TypeGoal && !grid.RawGrid.ValidTile(t)
						})
						if invalidGoals.Len() > 0 {
							continue
						}

						// all goals are valid!! send it over
						finalSolutionsCh <- solutionSet
					}
				}
				finalSolutionsWg.Done()
			}()
		}
	}

	go func() {
		finalSolutionsWg.Wait()
		close(finalSolutionsCh)
	}()

	return finalSolutionsCh
}

func findNewSolutionSets(
	baseGrid GridSolver,
	newSolution goalSolution,
	pairingsToUpdate [][][2]gs.Tile,
	currentSolutions map[[2]gs.Tile][]goalSolution,
) []gs.TileSet {
	newPair := [2]gs.Tile{newSolution.start, newSolution.end}
	var newSolutions []gs.TileSet

pairingsToUpdateLoop:
	for _, pairing := range pairingsToUpdate {
		solutions := make(map[[2]gs.Tile][]goalSolution, len(pairing)-1)
		for _, pair := range pairing {
			if pair != newPair {
				// if this pair does not have any solutions yet,
				// we do not care about this pairing
				if len(currentSolutions[pair]) == 0 {
					continue pairingsToUpdateLoop
				}
				solutions[pair] = currentSolutions[pair]
			} else {
				// we do not care about previous solutions found for this pair,
				// only add the new solution.
				solutions[pair] = []goalSolution{newSolution}
			}
		}

		forEachSolutionSet(solutions, func(solutionSet []goalSolution) {
			var finalSolution gs.TileSet
			for _, solution := range solutionSet {
				finalSolution.Merge(solution.path)
			}
			newSolutions = append(newSolutions, finalSolution)
		})
	}

	return newSolutions
}

func forEachSolutionSet(pairsToSolutions map[[2]gs.Tile][]goalSolution, forEach func([]goalSolution)) {
	for pair, solutions := range pairsToSolutions {

		// base case: call forEach on all solutions in pairsToSolutions
		if len(pairsToSolutions) == 1 {
			for _, solution := range solutions {
				forEach([]goalSolution{solution})
			}
			return
		}

		// recursive case: call recursively on all pairs except this one, then append each solution from
		// this pair onto all solutionSets from the recursive call
		remainingSolutions := make(map[[2]gs.Tile][]goalSolution, len(pairsToSolutions)-1)
		for pair2, solutions2 := range pairsToSolutions {
			if pairCompare(pair, pair2) < 0 {
				remainingSolutions[pair2] = solutions2
			}
		}

		forEachSolutionSet(remainingSolutions, func(solSet []goalSolution) {
			for _, solution := range solutions {
				newSolSet := make([]goalSolution, len(solSet)+1)
				newSolSet[0] = solution
				copy(newSolSet[1:], solSet)
				if anyIntersections(newSolSet) {
					return
				}
				forEach(newSolSet)
			}
		})
	}
}

func pairCompare(p1, p2 [2]gs.Tile) int8 {
	if p1 == p2 {
		return 0
	}

	if p1[0].X < p2[0].X {
		return -1
	}
	if p1[0].Y < p2[0].Y {
		return -1
	}

	if p1[1].X < p2[1].X {
		return -1
	}
	if p1[1].Y < p2[1].Y {
		return -1
	}

	return 1
}

func anyIntersections(solSet []goalSolution) bool {
	var allTiles gs.TileSet

	for _, solution := range solSet {
		for _, tile := range solution.path.Slice() {
			if allTiles.Has(tile) {
				return true
			}
			allTiles.Add(tile)
		}
	}

	return false
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
