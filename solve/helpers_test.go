package solve_test

import (
	"strings"
	"testing"

	gs "github.com/deanveloper/gridspech-go"
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

func tileSetFromString(grid gs.Grid, str string) gs.TileSet {
	lines := strings.Split(strings.Trim(str, "{}"), "|")
	var ts gs.TileSet
	for i, line := range lines {
		y := len(lines) - i - 1
		x := -1
		for {
			index := strings.IndexAny(line[x+1:], "012345")
			if index < 0 {
				break
			}
			x = x + 1 + index
			tileWithColor := *grid.TileAt(x, y)
			tileWithColor.Data.Color = gs.TileColor(line[x] - '0')
			ts.Add(tileWithColor)
		}
	}
	return ts
}

func commaSeparatedSlice(slice []gs.TileSet) string {
	var asStr []string
	for _, v := range slice {
		asStr = append(asStr, v.String())
	}
	return strings.Join(asStr, ",")
}

func testUnorderedTilesetSliceEq(t *testing.T, expected, actual []gs.TileSet) {
	t.Helper()

	for i1 := range expected {
		var found bool
		for i2 := range actual {
			if expected[i1].Eq(actual[i2]) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected to find tileset %v", expected[i1])
		}
	}

	for i1 := range actual {
		var found bool
		for i2 := range expected {
			if actual[i1].Eq(expected[i2]) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("incorrect solution %v", actual[i1])
		}
	}
}
