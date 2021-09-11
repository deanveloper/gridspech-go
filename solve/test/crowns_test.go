package solve_test

import (
	"testing"

	gs "github.com/deanveloper/gridspech-go"
	"github.com/deanveloper/gridspech-go/solve"
)

func testSolveCrownsAbstract(t *testing.T, level string, expectedSolutionsStrings []string, maxColors int) {
	t.Helper()

	solver := solve.NewGridSolver(gs.MakeGridFromString(level, maxColors))

	var expectedSolutions []gs.TileSet
	var actualSolutions []gs.TileSet

	for _, solutionString := range expectedSolutionsStrings {
		expectedSolutions = append(expectedSolutions, tileSetFromString(solver.Grid, solutionString))
	}

	solutionsCh := solver.SolveCrowns()
	for solution := range solutionsCh {
		actualSolutions = append(actualSolutions, solution)
	}

	testUnorderedTilesetSliceEq(t, expectedSolutions, actualSolutions)
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
