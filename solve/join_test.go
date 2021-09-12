package solve_test

import (
	"testing"

	"github.com/deanveloper/gridspech-go/solve"
)

func testSolveJoinsAbstract(t *testing.T, level string, expectedSolutionsStrings []string, maxColors int) {
	t.Helper()

	testSolveAbstract(t, level, expectedSolutionsStrings, maxColors, solve.GridSolver.SolveJoins)
}

func TestSolveJoins_basic(t *testing.T) {
	const level = `
	0j1  0  0  0j1
	`
	solutions := []string{
		"0000",
		"1111",
		"2222",
	}

	testSolveJoinsAbstract(t, level, solutions, 3)
}
