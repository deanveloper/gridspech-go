package main

import (
	"fmt"

	"github.com/deanveloper/gridspech-go"
	"github.com/deanveloper/gridspech-go/solve"
)

func main() {
	const levelAaa = `
[   ] [   ] [   ] [   ] [   ] [   ] [   ] [   ] 
[   ] [gA/] [   ] [   ] [gA/] [   ] [   ] [   ] 
[   ] [   ] [   ] [   ] [   ] [   ] [   ] [   ] 
[   ] [   ] [   ] [   ] [   ] [ A/] [   ] [   ] 
[   ] [   ] [   ] [   ] [   ] [g  ] [   ] [   ] 
[   ] [   ] [   ] [   ] [   ] [   ] [   ] [   ] 
[   ] [   ] [   ] [   ] [   ] [   ] [   ] [   ] 
[   ] [   ] [gA/] [   ] [   ] [   ] [   ] [   ] 
`
	grid := gridspech.MakeGridFromString(levelAaa)
	fmt.Println(grid)
	ch := solve.Goals(solve.Grid{Grid: grid}, 2)
	for solvedGrid := range ch {
		fmt.Println(solvedGrid)
		fmt.Println("=============")
	}
}
