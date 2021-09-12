package solve_test

import (
	"testing"

	gs "github.com/deanveloper/gridspech-go"
	"github.com/deanveloper/gridspech-go/solve"
)

func testSolveGoalsAbstract(t *testing.T, level string, maxColors int, expectedSolutionStrings []string) {
	t.Helper()

	solver := solve.NewGridSolver(gs.MakeGridFromString(level, 2))

	var expectedSolutions []gs.TileSet
	var actualSolutions []gs.TileSet

	for _, solutionString := range expectedSolutionStrings {
		expectedSolutions = append(expectedSolutions, tileSetFromString(solver.Grid, solutionString))
	}

	ch := solver.SolveGoals()
	for solution := range ch {
		actualSolutions = append(actualSolutions, solution)
	}

	testUnorderedTilesetSliceEq(t, expectedSolutions, actualSolutions)
}

func TestGoals_levelB1(t *testing.T) {
	const level = `
	0    0    0    0    0    0
	1/   0    0    0    0/   0
	1/e  0e   _    _    0e   0e
	`

	solutions := []string{
		"111111|100001|10  01",
	}
	testSolveGoalsAbstract(t, level, 2, solutions)
}

func TestGoals_levelB6(t *testing.T) {
	const level = `
	0e  0e  0   0
	0e  0   0   0
	0   0   0   0e
	0   0   0e  0e
	`

	solutions := []string{
		"001 |1101|0110| 010",
		"110 |0010|1001| 101",
		"010 |0110|1011| 100",
		"101 |1001|0100| 011",
	}
	testSolveGoalsAbstract(t, level, 2, solutions)
}

// test arrows with goals
func TestGoals_levelG3(t *testing.T) {
	const level = `
	_    0    0    _    _  
	_    0e   0^v  0e   0  
	0    0<>  _    0<>  0  
	0    0e   0^v  0e   _  
	_    _    0    0    _  
	`
	solutions := []string{
		" 00  | 1101|10 01|1011 |  00 ",
		" 01  | 1001|01 10|1001 |  10 ",
		" 11  | 0010|01 10|0100 |  11 ",
		" 10  | 0110|10 01|0110 |  01 ",
	}

	testSolveGoalsAbstract(t, level, 2, solutions)
}
