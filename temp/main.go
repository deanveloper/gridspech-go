package main

import (
	"fmt"

	"github.com/deanveloper/gridspech-go"
)

// func testNewC7() {
// 	const level = `
// [    ] [    ] [    ] [    ] [    ] [    ] [----]
// [g   ] [  / ] [    ] [  / ] [    ] [    ] [    ]
// [    ] [cA/ ] [    ] [g   ] [    ] [cA/ ] [    ]
// [    ] [    ] [    ] [    ] [ A/ ] [    ] [    ]
// [----] [    ] [g   ] [g   ] [    ] [    ] [    ]
// `
// 	grid := gridspech.MakeGridFromString(level, 2)
// 	solver := solve.NewGridSolver(grid)
// 	goals := solve.Goals(solver)
// 	for goalSolution := range goals {
// 		gridClone := solver.Grid.Clone()
// 		gridClone.ApplyTileSet(goalSolution)
// 		fmt.Println(gridClone)
// 		fmt.Println("")
// 	}
// }

// func testE4() {
// 	const level = `
// [2   ] [2   ] [2   ] [2   ] [2   ]
// [2   ] [2   ] [2   ] [2   ] [2   ]
// [2   ] [2   ] [    ] [g   ] [2   ]
// [2   ] [2   ] [g   ] [2   ] [2   ]
// [2   ] [2   ] [2   ] [2   ] [2   ]
// `
// 	grid := gridspech.MakeGridFromString(level, 2)
// 	solver := solve.NewGridSolver(grid)
// 	goals := solve.Goals(solver)
// 	for goalSolution := range goals {
// 		newSolver := solver.Clone()
// 		newSolver.Grid.ApplyTileSet(goalSolution)
// 		newSolver.UnknownTiles.RemoveAll(goalSolution.ToTileCoordSet())
// 		dots := solve.Dots(newSolver)
// 		for dotSolution := range dots {
// 			finalGrid := newSolver.Grid.Clone()
// 			fmt.Println(dotSolution.MultiLineString())
// 			finalGrid.ApplyTileSet(dotSolution)
// 			if finalGrid.Valid() {
// 				fmt.Println(finalGrid)
// 				fmt.Println("")
// 			}
// 		}
// 	}
// }

// func testNewE7() {
// 	const level = `
// [    ] [    ] [2   ] [    ] [    ]
// [    ] [    ] [3   ] [    ] [g   ]
// [    ] [    ] [gA/ ] [    ] [    ]
// [    ] [    ] [3   ] [    ] [    ]
// [    ] [    ] [2   ] [    ] [    ]
// [1   ] [    ] [c   ] [    ] [    ]
// `
// 	grid := gridspech.MakeGridFromString(level, 2)
// 	solver := solve.NewGridSolver(grid)
// 	goals := solve.Goals(solver)
// 	for goalSolution := range goals {
// 		newSolver := solver.Clone()
// 		newSolver.Grid.ApplyTileSet(goalSolution)
// 		dots := solve.Dots(newSolver)
// 		for dotSolution := range dots {
// 			finalGrid := newSolver.Grid.Clone()
// 			finalGrid.ApplyTileSet(dotSolution)
// 			if finalGrid.Valid() {
// 				fmt.Println(finalGrid)
// 				fmt.Println()
// 			}
// 		}
// 	}
// }

// func testNewE8() {
// 	const level = `
// [2   ] [2   ] [2   ] [2   ] [2   ] [2   ] [2   ]
// [2   ] [2   ] [3   ] [2   ] [2   ] [2   ] [2   ]
// [2   ] [2   ] [2   ] [2   ] [2   ] [2   ] [2   ]
// [2   ] [2   ] [2   ] [----] [2   ] [2   ] [2   ]
// [2   ] [2   ] [2   ] [2   ] [2   ] [3   ] [2   ]
// [2   ] [2   ] [2   ] [2   ] [2   ] [2   ] [2   ]
// [2   ] [2   ] [2   ] [2   ] [2   ] [2   ] [2   ]
// `

// 	grid := gridspech.MakeGridFromString(level, 2)
// 	solver := solve.NewGridSolver(grid)
// 	dots := solve.Dots(solver)
// 	for dotSolution := range dots {
// 		gridClone := solver.Grid.Clone()
// 		gridClone.ApplyTileSet(dotSolution)
// 		fmt.Println(gridClone)
// 		fmt.Println("")
// 	}
// }

// func testB1() {

// 	const level = `
// [    ] [    ] [    ] [    ] [    ] [    ]
// [ A/ ] [    ] [    ] [    ] [  / ] [    ]
// [gA/ ] [g   ] [----] [----] [g   ] [g   ]
// `
// 	grid := gridspech.MakeGridFromString(level, 2)
// 	solver := solve.NewGridSolver(grid)
// 	for solution := range solver.SolvePath(gridspech.TileCoord{X: 1, Y: 0}, gridspech.TileCoord{X: 4, Y: 0}, gridspech.ColorNone) {
// 		fmt.Println(solution)
// 	}
// }

// func testDotsLevel() {
// 	const level = `
// [1   ] [1   ] [2   ] [1   ] [1   ]
// [1   ] [1   ] [2   ] [1   ] [1   ]
// [1   ] [1   ] [    ] [1   ] [1   ]
// [    ] [2   ] [1   ] [2   ] [    ]
// [2   ] [1   ] [2   ] [1   ] [2   ]
// `
// 	grid := gridspech.MakeGridFromString(level, 2)
// 	solver := solve.NewGridSolver(grid)
// 	for solution := range solve.Dots(solver) {
// 		newGrid := grid.Clone()
// 		newGrid.ApplyTileSet(solution)
// 		fmt.Println(newGrid)
// 		fmt.Println()
// 	}
// }

// func pathTest() {

// 	const levelAaa = `
// [    ] [    ] [    ] [    ] [    ] [    ]
// [    ] [gA/ ] [    ] [    ] [gA/ ] [    ]
// [    ] [ A/ ] [    ] [    ] [ A/ ] [    ]
// [    ] [ A/ ] [    ] [  / ] [    ] [ A/ ]
// [    ] [ A/ ] [    ] [  / ] [    ] [gA/ ]
// [    ] [ A/ ] [    ] [    ] [    ] [    ]
// [    ] [    ] [ A/ ] [    ] [    ] [    ]
// [    ] [    ] [gA/ ] [    ] [    ] [    ]
// `
// 	grid := gridspech.MakeGridFromString(levelAaa, 2)
// 	ch := solve.Goals(solve.NewGridSolver(grid))
// 	for solvedGrid := range ch {
// 		fmt.Println(solvedGrid.MultiLineString())
// 		fmt.Println("=============")
// 	}
// }

// func standardSerialization() {
// 	const lv = `
// 	0    0    0    0    0    0
// 	0    0    0    0    0    0
// 	0    0    0e   0    0/   0
// 	0    0    0    0    0e   0
// 	0    1/e  0e   1/   0    0
// 	`

// 	grid := gridspech.DeserializeGrid(lv, 2)
// 	ch := solve.Goals(solve.NewGridSolver(grid))
// 	for solvedGrid := range ch {
// 		fmt.Println(solvedGrid.MultiLineString())
// 		fmt.Println("=============")
// 	}
// }

func convertToStandard() {
	const level = `
[ A/ ] [    ] [    ] [ A/ ]
[g / ] [    ] [    ] [g   ]
[ A/ ] [    ] [    ] [ A/ ]
`
	grid := gridspech.MakeGridFromStringOld(level, 3)
	fmt.Println(grid.String())
}

func main() {
	convertToStandard()
}
