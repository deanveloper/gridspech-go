package solve_test

import (
	"testing"

	gs "github.com/deanveloper/gridspech-go"
	"github.com/deanveloper/gridspech-go/solve"
)

func testSolveDotsAbstract(t *testing.T, level string, maxColors int, solutions []string) {
	t.Helper()

	grid := solve.NewGridSolver(gs.MakeGridFromString(level, 2))

	actualSolutions := solve.Dots(grid)
	var actualSolutionsStrs []string
	for solution := range actualSolutions {
		solvedGrid := grid.Grid.Clone()
		solvedGrid.ApplyTileSet(solution)

		actualSolutionsStrs = append(actualSolutionsStrs, gridToSolutionString(t, solvedGrid))
	}

	testStringSlicesEq(t, solutions, actualSolutionsStrs)
}

func TestDots_levelBasic1(t *testing.T) {
	const level = `
	0    0    0  
	0    0m1  0  
	0    0    0  
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
	0    0    0  
	0    0m2  0  
	0    0    0  
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
	0m1  0m1  0m2  0m1  0m1
	0m1  0m1  0m2  0m1  0m1
	0m1  0m1  0    0m1  0m1
	0    0m2  0m1  0m2  0
	0m2  0m1  0m2  0m1  0m2
	`
	solutions := []string{
		" A A \n A A \n     \nA   A\n AAA ",
	}
	testSolveDotsAbstract(t, level, 2, solutions)

}
