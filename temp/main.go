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
	ch := solve.Goals(solve.NewGridSolver(grid))
	for solvedGrid := range ch {
		fmt.Println(solvedGrid.MultiLineString())
		fmt.Println("=============")
	}
}

func main() {
	standardSerialization()
}
