package solve_test

import (
	"testing"

	gs "github.com/deanveloper/gridspech-go"
	"github.com/deanveloper/gridspech-go/solve"
)

func testSolveDotsAbstract(t *testing.T, level string, maxColors int, expectedSolutionStrings []string) {
	t.Helper()

	solver := solve.NewGridSolver(gs.MakeGridFromString(level, 2))

	var expectedSolutions []gs.TileSet
	var actualSolutions []gs.TileSet

	for _, solutionString := range expectedSolutionStrings {
		expectedSolutions = append(expectedSolutions, tileSetFromString(solver.Grid, solutionString))
	}

	for solution := range solver.SolveDots() {
		actualSolutions = append(actualSolutions, solution)
	}

	testUnorderedTilesetSliceEq(t, expectedSolutions, actualSolutions)
}

func TestDots_levelBasic1(t *testing.T) {
	const level = `
	0    0    0  
	0    0m1  0  
	0    0    0  
	`
	solutions := []string{
		" 1 |0 0| 0 ",
		" 0 |1 0| 0 ",
		" 0 |0 1| 0 ",
		" 0 |0 0| 1 ",
	}
	testSolveDotsAbstract(t, level, 2, solutions)
}

func TestDots_levelBasic2(t *testing.T) {
	const level = `
	0    0    0  
	0    0m2  0  
	0    0    0  
	`
	solutions := []string{
		" 1 |1 0| 0 ",
		" 1 |0 1| 0 ",
		" 1 |0 0| 1 ",
		" 0 |1 1| 0 ",
		" 0 |1 0| 1 ",
		" 0 |0 1| 1 ",
	}
	testSolveDotsAbstract(t, level, 2, solutions)
}

func TestDots_levelE8(t *testing.T) {
	const level = `
	0m1  0m1  0m2  0m1  0m1
	0m1  0m1  0m2  0m1  0m1
	0m1  0m1  0    0m1  0m1
	0    0m2  0m1  0m2  0
	0m2  0m1  0m2  0m1  0m2
	`
	solutions := []string{
		"01010|01010|00000|10001|01110",
	}
	testSolveDotsAbstract(t, level, 2, solutions)

}
