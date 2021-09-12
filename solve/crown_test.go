package solve_test

import (
	"testing"

	"github.com/deanveloper/gridspech-go/solve"
)

func testSolveCrownsAbstract(t *testing.T, level string, expectedSolutionsStrings []string, maxColors int) {
	t.Helper()

	testSolveAbstract(t, level, expectedSolutionsStrings, maxColors, solve.GridSolver.SolveCrowns)
}

func TestSolveCrowns_basic1(t *testing.T) {
	const lvl = `
	0  0  0k
	`
	solutions := []string{
		"000", " 10", "100", " 20", "200",
		"111", " 01", " 21", "011", "211",
		"222", " 02", " 12", "022", "122",
	}

	testSolveCrownsAbstract(t, lvl, solutions, 3)
}
