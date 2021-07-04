package solve

import (
	"sync"

	"github.com/deanveloper/gridspech-go"
	gs "github.com/deanveloper/gridspech-go"
)

type solutionsChan struct {
	color gs.TileColor
	sols  <-chan gs.TileSet
}
type solutionForColor struct {
	color gs.TileColor
	sol   gs.TileSet
}

// Goals will return a channels of solutions for all the goal tiles in g
func Goals(g Grid, maxColors gs.TileColor) <-chan Grid {
	var goalTiles []Tile

	for _, col := range g.Tiles {
		for _, tile := range col {
			if tile.Type == gridspech.TypeGoal {
				goalTiles = append(goalTiles, tile)
			}
		}
	}

	var solutionsChans []solutionsChan
	for i1 := 0; i1 < len(goalTiles)-1; i1++ {
		for i2 := i1 + 1; i2 < len(goalTiles); i2++ {
			goal1, goal2 := goalTiles[i1], goalTiles[i2]

			for color := gs.TileColor(0); color < maxColors; color++ {
				paths := g.SolvePath(goal1, goal2, color)
				solutionsChans = append(solutionsChans, solutionsChan{color: color, sols: paths})
			}
		}
	}

	return solutionsToGrids(g, aggregateSolutions(solutionsChans))
}

func solutionsToGrids(baseGrid Grid, sols <-chan solutionForColor) <-chan Grid {
	ch := make(chan Grid)
	go func() {
		for sol := range sols {
			newGrid := Grid{Grid: baseGrid.Clone()}
			for _, tile := range sol.sol.Slice() {
				newGrid.Tiles[tile.X][tile.Y].Color = sol.color
			}
			ch <- newGrid
		}
		close(ch)
	}()
	return ch
}

func aggregateSolutions(solutions []solutionsChan) <-chan solutionForColor {
	ch := make(chan solutionForColor)

	var wg sync.WaitGroup
	wg.Add(len(solutions))

	mergeIn := func(solsCh solutionsChan) {
		for sol := range solsCh.sols {
			ch <- solutionForColor{
				color: solsCh.color,
				sol:   sol,
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
