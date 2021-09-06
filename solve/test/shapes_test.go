package solve_test

import (
	"testing"

	gs "github.com/deanveloper/gridspech-go"
	"github.com/deanveloper/gridspech-go/solve"
)

func testShapesIterAbstract(t *testing.T, level string, start gs.TileCoord, color gs.TileColor, expectedSolutionsStrings []string) {
	t.Helper()

	solver := solve.NewGridSolver(gs.MakeGridFromString(level, 2))

	var expectedSolutions []gs.TileSet
	var actualSolutions []gs.TileSet

	for _, solutionString := range expectedSolutionsStrings {
		expectedSolutions = append(expectedSolutions, tileSetFromString(solver.Grid, solutionString))
	}

	solutionsCh, prune := solver.ShapesIter(start, color)
	for solution := range solutionsCh {
		actualSolutions = append(actualSolutions, solution)
		prune <- false
	}

	testUnorderedTilesetSliceEq(t, expectedSolutions, actualSolutions)
}

func TestShapesIter_small(t *testing.T) {
	const level = `
	0  0
	0  0
	`

	solutions := []string{
		"  | 1", " 1| 1", "11| 1",
		"  |11", " 1|11", "11|11", "1 |11",
	}

	testShapesIterAbstract(t, level, gs.TileCoord{X: 1, Y: 0}, 1, solutions)
}

func TestShapesIter_large(t *testing.T) {
	const level = `
	0  0  0  
	0  0  0  
	`

	solutions := []string{
		"   |  1", "  1|  1", " 11|  1", "111|  1",
		"   | 11", "  1| 11", " 11| 11", "111| 11", "11 | 11", " 1 | 11",
		"   |111", "  1|111", " 11|111", "111|111", "11 |111", " 1 |111", "1  |111", "1 1|111",
		"111|1 1",
	}

	testShapesIterAbstract(t, level, gs.TileCoord{X: 2, Y: 0}, 1, solutions)
}

func TestShapesIter_withHole(t *testing.T) {
	const level = `
	_  0  0  
	0  0  0  
	`

	solutions := []string{
		"   |  1", "  1|  1", " 11|  1",
		"   | 11", "  1| 11", " 11| 11", " 1 | 11",
		"   |111", "  1|111", " 11|111", " 1 |111",
	}

	testShapesIterAbstract(t, level, gs.TileCoord{X: 2, Y: 0}, 1, solutions)
}

func TestShapesIter_noTraverseKnownDifferent(t *testing.T) {
	const level = `
	0   0   0  
	0   0/  0/  
	0/  0   0  
	`

	solutions := []string{"  1", " 11"}

	testShapesIterAbstract(t, level, gs.TileCoord{X: 2, Y: 0}, 1, solutions)
}

func TestShapesIter_traverseKnownSame(t *testing.T) {
	const level = `
	0   1/  0   
	`

	solutions := []string{"  1", " 11", "111"}

	testShapesIterAbstract(t, level, gs.TileCoord{X: 2, Y: 0}, 1, solutions)
}

func TestShapesIter_decoratedBorderSmall(t *testing.T) {

	const level = `
	0  0
	0  0
	`

	solutions := []string{
		" 1|10",
		"10|10", "11|00",
		"10|00", "00|10", "01|00",
		"00|00",
	}

	testShapesIterAbstract(t, level, gs.TileCoord{X: 1, Y: 0}, 0, solutions)
}

func TestShapesIter_decoratedBorderLarge(t *testing.T) {
	const level = `
	0  0  0  
	0  0  0  
	`

	solutions := []string{
		"  1| 10", " 10| 10", "100| 10", "000|110",
		" 11|100", " 10|100", "100|100", "000|100", "001|100", "101|100",
		"111|000", "110|000", "100|000", "000|000", "001|000", "101|000", "011|000", "010|000",
		"000|010",
	}

	testShapesIterAbstract(t, level, gs.TileCoord{X: 2, Y: 0}, 0, solutions)
}

func TestShapesIter_noDuplicates(t *testing.T) {
	const level = `
	0  0  0  
	0  0  0  
	0  0  0  
	`
	solver := solve.NewGridSolver(gs.MakeGridFromString(level, 2))

	var solutions []gs.TileSet

	solutionsCh, pruneCh := solver.ShapesIter(gs.TileCoord{X: 1, Y: 1}, 1)
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
