package solve_test

import (
	"testing"

	gs "github.com/deanveloper/gridspech-go"
	"github.com/deanveloper/gridspech-go/solve"
)

func testSolveDotsAbstract(t *testing.T, level string, maxColors int, solutions []string) {
	t.Helper()

	grid := solve.NewGridSolver(gs.MakeGridFromString(level))

	actualSolutions := solve.Dots(grid, maxColors)
	var actualSolutionsStrs []string
	for solution := range actualSolutions {
		solvedGrid := grid.Grid()
		solvedGrid.ApplyTileSet(solution)

		actualSolutionsStrs = append(actualSolutionsStrs, gridToSolutionString(t, solvedGrid))
	}

	testStringSlicesEq(t, solutions, actualSolutionsStrs)
}

func TestDots_levelBasic1(t *testing.T) {
	const level = `
[    ] [    ] [    ] 
[    ] [1   ] [    ] 
[    ] [    ] [    ] 
`
	solutions := []string{
		" A \n   \n   ",
		"   \nA  \n   ",
		"   \n  A\n   ",
		"   \n   \n A ",
	}
	testSolveDotsAbstract(t, level, 2, solutions)
}

func TestDots_levelBasic2(t *testing.T) {
	const level = `
[    ] [    ] [    ] 
[    ] [2   ] [    ] 
[    ] [    ] [    ] 
`
	solutions := []string{
		" A \nA  \n   ",
		" A \n  A\n   ",
		" A \n   \n A ",
		"   \nA A\n   ",
		"   \nA  \n A ",
		"   \n  A\n A ",
	}
	testSolveDotsAbstract(t, level, 2, solutions)
}

func TestDots_levelE8(t *testing.T) {
	const level = `
[1   ] [1   ] [2   ] [1   ] [1   ]
[1   ] [1   ] [2   ] [1   ] [1   ]
[1   ] [1   ] [    ] [1   ] [1   ]
[    ] [2   ] [1   ] [2   ] [    ]
[2   ] [1   ] [2   ] [1   ] [2   ]
`
	solutions := []string{
		" x x \n x x \n     \nx   x\n xxx ",
	}
	testSolveDotsAbstract(t, level, 2, solutions)

}

func BenchmarkDots_levelBasic2(b *testing.B) {
	const level = `
[    ] [    ] [    ] 
[    ] [2   ] [    ] 
[    ] [    ] [    ] 
`
	for x := 0; x < b.N; x++ {
		grid := solve.NewGridSolver(gs.MakeGridFromString(level))
		solve.Dots(grid, 2)
	}
}

func BenchmarkDots_levelE8(b *testing.B) {
	const level = `
[1   ] [1   ] [2   ] [1   ] [1   ]
[1   ] [1   ] [2   ] [1   ] [1   ]
[1   ] [1   ] [    ] [1   ] [1   ]
[    ] [2   ] [1   ] [2   ] [    ]
[2   ] [1   ] [2   ] [1   ] [2   ]
`
	for x := 0; x < b.N; x++ {
		grid := solve.NewGridSolver(gs.MakeGridFromString(level))
		solve.Dots(grid, 2)
	}
}
