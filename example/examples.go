package example

import (
	"fmt"
	"sort"

	"github.com/deanveloper/gridspech-go"
	"github.com/deanveloper/gridspech-go/solve"
)

// Examples of all A levels in gridspech
const (
	LevelA1 = `1/e  0    0    0e`
	LevelA2 = `
	1/e  _    0    0e 
	0    0    0    _  
	`
	LevelA3 = `
	_    0    0    0    _  
	1/e  0    0/   0    0e 
	_    0    0    0    _  
	`
	LevelA4 = `
	1/e  0    0e 
	0    0    1/ 
	`
	LevelA5 = `
	0    1/   0    0  
	0    0    1/   0  
	1/e  0    0/   0e 
	`
	LevelA6 = `
	0    0    0    0    1/   0    0    0  
	0    0    1/   0    1/   0    1/   0  
	1/e  0    0    0    1/   0    0    0e 
	`
	LevelA7 = `
	_    0    1/   0    _  
	_    0    1/   0    _  
	1/e  0    0    0    0e 
	_    0    0    0    _  
	`
	LevelA8 = `
	0    0    1/   0    0  
	0    1/   0    1/   0  
	1/e  0    0    0    1/e
	0    1/   0    1/   0  
	0    0    1/   0    0  
	`
	LevelA9 = `
	0    1/   0    0    1/   0    0  
	1/e  0    0    1/   0    0    0e 
	0    0    1/   0    1/   0    0  
	`
)

// FindSolutionsForALevels finds the solution to all A levels in gridspech
func FindSolutionsForALevels() {
	levelMap := map[string]string{
		"A1": LevelA1,
		"A2": LevelA2,
		"A3": LevelA3,
		"A4": LevelA4,
		"A5": LevelA5,
		"A6": LevelA6,
		"A7": LevelA7,
		"A8": LevelA8,
		"A9": LevelA9,
	}
	levels := make([]string, 0, len(levelMap))
	for k := range levelMap {
		levels = append(levels, k)
	}
	sort.Strings(levels)

	for _, level := range levels {
		levelStr := levelMap[level]
		solution := FindSolution(gridspech.MakeGridFromString(levelStr, 2))
		fmt.Println("solution for level " + level + ":")
		fmt.Println(solution)
	}
}

// FindSolution returns a new grid as a solution to `grid`
func FindSolution(grid gridspech.Grid) gridspech.Grid {
	ch := solve.Goals(solve.NewGridSolver(grid))
	newGrid := grid.Clone()
	firstSolution := <-ch
	newGrid.ApplyTileSet(firstSolution)

	return newGrid
}
