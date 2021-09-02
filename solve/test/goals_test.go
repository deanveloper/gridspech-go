package solve_test

import (
	"strings"
	"testing"

	gs "github.com/deanveloper/gridspech-go"
	"github.com/deanveloper/gridspech-go/solve"
)

func gridToSolutionString(t *testing.T, g gs.Grid) string {
	t.Helper()

	byteArray := make([]byte, g.Height()*(g.Width()+1))
	for x, col := range g.Tiles {
		for y, tile := range col {
			index := x + (g.Width()+1)*(g.Height()-y-1)
			if tile.Data.Color == gs.ColorNone {
				byteArray[index] = ' '
			} else {
				byteArray[index] = byte(g.TileAt(x, y).Data.Color) + 'A' - 1
			}

			if x == g.Width()-1 {
				byteArray[index+1] = '\n'
			}
		}
	}

	return string(byteArray)
}

func testStringSlicesEq(t *testing.T, expected, actual []string) {
	t.Helper()

	for i1 := range expected {
		var found bool
		for i2 := range actual {
			if strings.TrimSpace(expected[i1]) == strings.TrimSpace(actual[i2]) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected to find solution, but did not:\n%v", expected[i1])
		}
	}

	for i1 := range actual {
		var found bool
		for i2 := range expected {
			if strings.TrimSpace(actual[i1]) == strings.TrimSpace(expected[i2]) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("solution not expected:\n%v", actual[i1])
		}
	}
}

func testSolveGoalsAbstract(t *testing.T, level string, maxColors int, solutions []string) {
	t.Helper()

	grid := solve.NewGridSolver(gs.MakeGridFromString(level, 2))

	ch := solve.Goals(grid)
	var actualSolutions []string
	for solution := range ch {
		solvedGrid := grid.Grid.Clone()
		solvedGrid.ApplyTileSet(solution)

		actualSolutions = append(actualSolutions, gridToSolutionString(t, solvedGrid))
	}

	testStringSlicesEq(t, solutions, actualSolutions)
}

func TestGoals_levelB1(t *testing.T) {
	const level = `
	0    0    0    0    0    0
	1/   0    0    0    0/   0
	1/e  0e   _    _    0e   0e
	`

	solutions := []string{
		"AAAAAA\nA    A\nA    A",
	}
	testSolveGoalsAbstract(t, level, 2, solutions)
}

func TestGoals_levelB6(t *testing.T) {
	const level = `
	0e  0e  0   0
	0e  0   0   0
	0   0   0   0e
	0   0   0e  0e
	`

	solutions := []string{
		"  A \nAA A\n AA \n  A ",
		" A  \n AA \nA AA\n A  ",
		"A A \nA  A\n A  \n  AA",
		"AA  \n  A \nA  A\n A A",
	}
	testSolveGoalsAbstract(t, level, 2, solutions)
}
