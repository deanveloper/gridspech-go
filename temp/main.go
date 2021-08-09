package main

import (
	"fmt"

	"github.com/deanveloper/gridspech-go"
	"github.com/deanveloper/gridspech-go/solve"
)

func main() {
	const levelAaa = `
[    ] [    ] [    ] [    ] [    ] [    ]
[    ] [gA/ ] [    ] [    ] [gA/ ] [    ]
[    ] [ A/ ] [    ] [    ] [ A/ ] [    ]
[    ] [ A/ ] [    ] [  / ] [    ] [ A/ ]
[    ] [ A/ ] [    ] [  / ] [    ] [gA/ ]
[    ] [ A/ ] [    ] [    ] [    ] [    ]
[    ] [    ] [ A/ ] [    ] [    ] [    ]
[    ] [    ] [gA/ ] [    ] [    ] [    ]
`
	grid := gridspech.MakeGridFromString(levelAaa)
	ch := solve.Goals(solve.NewGridSolver(grid), 2)
	for solvedGrid := range ch {
		fmt.Println(solvedGrid.MultiLineString())
		fmt.Println("=============")
	}
}
