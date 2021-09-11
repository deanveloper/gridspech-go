package main

import (
	"fmt"

	"github.com/deanveloper/gridspech-go"
	"github.com/deanveloper/gridspech-go/solve"
)

func standardSerialization() {
	const lv = `
	0    0    0    0    0    0
	0    0    0    0    0    0
	0    0    0e   0    0/   0
	0    0    0    0    0e   0
	0    1/e  0e   1/   0    0
	`

	grid := gridspech.MakeGridFromString(lv, 2)
	ch := solve.NewGridSolver(grid).SolveGoals()
	for solvedGrid := range ch {
		fmt.Println(solvedGrid.MultiLineString())
		fmt.Println("=============")
	}
}

func tryE8() {
	const lvl = `
	0m2  0m2  0m2  0m2  0m2  0m2  0m2  
	0m2  0m2  0m3  0m2  0m2  0m2  0m2  
	0m2  0m2  0m2  0m2  0m2  0m2  0m2  
	0m2  0m2  0m2  _    0m2  0m2  0m2  
	0m2  0m2  0m2  0m2  0m2  0m3  0m2  
	0m2  0m2  0m2  0m2  0m2  0m2  0m2  
	0m2  0m2  0m2  0m2  0m2  0m2  0m2  
	`
	grid := gridspech.MakeGridFromString(lvl, 2)
	ch := solve.NewGridSolver(grid).SolveDots()
	for solvedGrid := range ch {
		fmt.Println(solvedGrid.MultiLineString())
		fmt.Println("=============")
	}
}

func tryTest() {
	const lvl = `0   0   0k`
	grid := gridspech.MakeGridFromString(lvl, 3)
	ch := solve.NewGridSolver(grid).SolveCrowns()
	for solvedGrid := range ch {
		fmt.Println(solvedGrid.MultiLineString())
		fmt.Println("=============")
	}
}

func main() {
	tryTest()
}
