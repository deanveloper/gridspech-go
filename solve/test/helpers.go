package solve_test

import (
	"testing"

	"github.com/deanveloper/gridspech-go/solve"
)

func testSolveAbstract(t *testing.T, level string, expectedSolutionsStrings []string, maxColors int, f func(solve.GridSolver) <-chan gs.TileSet) {
	t.Helper()

	solver := solve.NewGridSolver(gs.MakeGridFromString(level, maxColors))

	var expectedSolutions []gs.TileSet
	var actualSolutions []gs.TileSet

	for _, solutionString := range expectedSolutionsStrings {
		expectedSolutions = append(expectedSolutions, tileSetFromString(solver.Grid, solutionString))
	}

	solutionsCh := f(solver)
	for solution := range solutionsCh {
		actualSolutions = append(actualSolutions, solution)
	}

	testUnorderedTilesetSliceEq(t, expectedSolutions, actualSolutions)
}

func testSolveCrownsAbstract(t *testing.T, level string, expectedSolutionsStrings []string, maxColors int) {
	t.Helper()

	testSolveAbstract(t, level, expectedSolutionsStrings, maxColors, solve.GridSolver.SolveCrowns)
}

func testSolveJoinsAbstract(t *testing.T, level string, expectedSolutionsStrings []string, maxColors int) {
	t.Helper()

	testSolveAbstract(t, level, expectedSolutionsStrings, maxColors, solve.GridSolver.SolveJoins)
}
