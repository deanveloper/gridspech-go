package solve_test

import (
	"testing"

	gs "github.com/deanveloper/gridspech-go"
	"github.com/deanveloper/gridspech-go/solve"
)

func testSolveDotsAbstract(t *testing.T, level string, maxColors int, expectedSolutionStrings []string) {
	t.Helper()

	testSolveAbstract(t, level, expectedSolutionStrings, maxColors, solve.GridSolver.SolveDots)
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

// test arrows with dots
func TestDots_levelG4(t *testing.T) {
	const level = `
	0m3<^v>  _        0<^v>    _        _        _        0<^v>    _      
	0<^v>    0m2^v    0m1^v    0m3^v    0^v      0m2^v    0m3^v    0m1<^v>
	`
	solutions := []string{"  0   1 |10110100"}
	testSolveDotsAbstract(t, level, 2, solutions)
}

func TestDots_noDuplicates(t *testing.T) {
	const level = `
	0m1  0m1  0m2  0m1  0m1
	0m1  0m1  0m2  0m1  0m1
	0m1  0m1  0    0m1  0m1
	0    0m2  0m1  0m2  0
	0m2  0m1  0m2  0m1  0m2
	`
	var prevSolutions []gs.TileSet
	for solution := range solve.NewGridSolver(gs.MakeGridFromString(level, 2)).SolveDots() {
		for _, prev := range prevSolutions {
			if prev.Eq(solution) {
				t.Fatalf("got duplicate: %v", prev)
			}
			prevSolutions = append(prevSolutions, solution)
		}
	}
}
