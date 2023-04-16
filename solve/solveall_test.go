package solve_test

import (
	"testing"

	"github.com/deanveloper/gridspech-go/solve"
)

func testSolveAllTilesAbstract(t *testing.T, level string, expectedSolutionsStrings []string, maxColors int) {
	t.Helper()

	testSolveAbstract(t, level, expectedSolutionsStrings, maxColors, solve.GridSolver.SolveAllTiles)
}

func TestSolveAllTiles_levelF10(t *testing.T) {
	const level = `
	0    0    0    0    0    0  
	0    0e   0k   0    0    0  
	0    0    0k   0    0    0  
	0    0    0k   0    0    0  
	0    0    0k   0e   0    0  
	0j1  0e   0    0    0j1  0e 
	`
	solutions := []string{
		"111111|100001|111100|100010|101010|101110",
		"000000|011110|000011|011101|010101|010001",
	}

	testSolveAllTilesAbstract(t, level, solutions, 2)
}
