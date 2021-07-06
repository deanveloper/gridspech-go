package solve_test

import (
	"strings"
	"testing"

	gs "github.com/deanveloper/gridspech-go"
	"github.com/deanveloper/gridspech-go/solve"
)

func gridToSolutionString(t *testing.T, g solve.Grid) string {
	t.Helper()

	byteArray := make([]byte, g.Height()*(g.Width()+1))
	for x, col := range g.Tiles {
		for y, tile := range col {
			index := x + (g.Height() - y - 1)
			if tile.Color == gs.ColorNone {
				byteArray[index] = ' '
			} else {
				byteArray[index] = byte(g.Tiles[x][y].Color) + 'A' - 1
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

func testSolveGoalsAbstract(t *testing.T, level string, maxColors gs.TileColor, solutions []string) {
	t.Helper()

	grid := solve.Grid{Grid: gs.MakeGridFromString(level)}

	ch := solve.Goals(grid, maxColors)
	var actualSolutions []string
	for solvedGrid := range ch {
		actualSolutions = append(actualSolutions, gridToSolutionString(t, solvedGrid))
	}

	testStringSlicesEq(t, solutions, actualSolutions)
}

func TestGoals_levelB1(t *testing.T) {
	const level = `
[   ] [   ] [   ] [   ] [   ] [   ] 
[ A/] [   ] [   ] [   ] [  /] [   ] 
[gA/] [g  ]             [g  ] [g  ] 
`

	solutions := []string{
		"AAAAAA\nA    A\nA    A",
	}
	testSolveGoalsAbstract(t, level, 2, solutions)
}
