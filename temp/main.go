package main

import (
	"fmt"
	"strings"

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
	ch := solve.NewGridSolver(grid).SolveEnds()
	for solvedGrid := range ch {
		fmt.Println(solvedGrid.MultiLineString())
		fmt.Println("=============")
	}
}

func flipRows() {
	const level = `
	0k    2     0     0     1   
	2     0     0     1k    0   
	0     0     0     1     0j1 
	0m2   1     0     0     1   
	0m1   0     0     0     0   
	0m2   0     0     0     0   
	0     _     0     1e    0   
	2j1   2/j1  2j1   1e    1   
	`

	split := strings.Split(level, "\n")
	for i := range split {
		split[i] = strings.TrimSpace(split[i])
	}
	for i := 0; i < len(split)/2; i++ {
		j := len(split) - i - 1
		split[i], split[j] = split[j], split[i]
	}
	fmt.Println(strings.Join(split[1:len(split)-1], "\n"))
}

func main() {
	flipRows()
}
