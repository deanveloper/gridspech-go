package solve_test

import (
	"testing"

	gs "github.com/deanveloper/gridspech-go"
	"github.com/deanveloper/gridspech-go/solve"
)

func testSolveShapesAbstract(t *testing.T, level string, start gs.TileCoord, color gs.TileColor, expectedSolutionsStrings []string) {
	t.Helper()

	solver := solve.NewGridSolver(gs.MakeGridFromString(level, 2))

	var expectedSolutions []gs.TileSet
	var actualSolutions []gs.TileSet

	for _, solutionString := range expectedSolutionsStrings {
		expectedSolutions = append(expectedSolutions, tileSetFromString(solver.Grid, solutionString))
	}

	solutionsCh, prune := solver.SolveShapes(start, color)
	for solution := range solutionsCh {
		actualSolutions = append(actualSolutions, solution)
		prune <- false
	}

	testUnorderedTilesetSliceEq(t, expectedSolutions, actualSolutions)
}

func TestSolveShapes_small(t *testing.T) {
	const level = `
	0  0
	0  0
	`

	solutions := []string{
		"  | x", " x| x", "xx| x",
		"  |xx", " x|xx", "xx|xx", "x |xx",
	}

	testSolveShapesAbstract(t, level, gs.TileCoord{X: 1, Y: 0}, 1, solutions)
}

func TestSolveShapes_large(t *testing.T) {
	const level = `
	0  0  0  
	0  0  0  
	`

	solutions := []string{
		"   |  x", "  x|  x", " xx|  x", "xxx|  x",
		"   | xx", "  x| xx", " xx| xx", "xxx| xx", "xx | xx", " x | xx",
		"   |xxx", "  x|xxx", " xx|xxx", "xxx|xxx", "xx |xxx", " x |xxx", "x  |xxx", "x x|xxx",
		"xxx|x x",
	}

	testSolveShapesAbstract(t, level, gs.TileCoord{X: 2, Y: 0}, 1, solutions)
}

func TestSolveShapes_withHole(t *testing.T) {
	const level = `
	_  0  0  
	0  0  0  
	`

	solutions := []string{
		"   |  x", "  x|  x", " xx|  x",
		"   | xx", "  x| xx", " xx| xx", " x | xx",
		"   |xxx", "  x|xxx", " xx|xxx", " x |xxx",
	}

	testSolveShapesAbstract(t, level, gs.TileCoord{X: 2, Y: 0}, 1, solutions)
}

func TestSolveShapes_noDuplicates(t *testing.T) {
	const level = `
	_  0  0  
	0  0  0  
	`
	solver := solve.NewGridSolver(gs.MakeGridFromString(level, 2))

	var solutions []gs.TileSet

	solutionsCh, pruneCh := solver.SolveShapes(gs.TileCoord{X: 2, Y: 0}, 1)
	for sol := range solutionsCh {
		for _, oldSol := range solutions {
			if sol.Eq(oldSol) {
				t.Errorf("Duplicate found: %v", oldSol)
				break
			}
		}
		solutions = append(solutions, sol)
		pruneCh <- false
	}
}
