package main

import (
	"fmt"
	"sort"

	"github.com/deanveloper/gridspech-go"
	"github.com/deanveloper/gridspech-go/solve"
)

// Examples of all A levels in gridspech
const (
	LevelA1 = `[gA/] [   ] [   ] [g  ]`
	LevelA2 = `
[gA/]       [   ] [g  ]
[   ] [   ] [   ]      
`
	LevelA3 = `
	[   ] [   ] [   ]      
[gA/] [   ] [  /] [   ] [g  ]
	[   ] [   ] [   ]      `
	LevelA4 = `
[gA/] [   ] [g  ]
[   ] [   ] [ A/]`
	LevelA5 = `
[   ] [ A/] [   ] [   ]
[   ] [   ] [ A/] [   ]
[gA/] [   ] [  /] [g  ]`
	LevelA6 = `
[   ] [   ] [   ] [   ] [ A/] [   ] [   ] [   ]
[   ] [   ] [ A/] [   ] [ A/] [   ] [ A/] [   ]
[gA/] [   ] [   ] [   ] [ A/] [   ] [   ] [g  ]`
	LevelA7 = `
      [   ] [ A/] [   ]       
      [   ] [ A/] [   ]       
[gA/] [   ] [   ] [   ] [g  ] 
      [   ] [   ] [   ]       `
	LevelA8 = `
[   ] [   ] [ A/] [   ] [   ] 
[   ] [ A/] [   ] [ A/] [   ] 
[gA/] [   ] [   ] [   ] [gA/] 
[   ] [ A/] [   ] [ A/] [   ] 
[   ] [   ] [ A/] [   ] [   ]`
	LevelA9 = `
[   ] [ A/] [   ] [   ] [ A/] [   ] [   ]
[gA/] [   ] [   ] [ A/] [   ] [   ] [g  ]
[   ] [   ] [ A/] [   ] [ A/] [   ] [   ]
`
)

func main() {
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
		solution := FindSolution(gridspech.MakeGridFromString(levelStr))
		fmt.Println("solution for level " + level + ":")
		fmt.Println(solution)
	}
}

// FindSolution returns a new grid as a solution to `grid`
func FindSolution(grid gridspech.Grid) gridspech.Grid {
	ch := solve.Goals(solve.NewGridSolver(grid), 2)
	newGrid := grid.Clone()
	firstSolution := <-ch
	newGrid.ApplyTileSet(firstSolution)

	return newGrid
}
