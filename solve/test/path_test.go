package solve_test

import (
	"strings"
	"testing"

	gs "github.com/deanveloper/gridspech-go"
	"github.com/deanveloper/gridspech-go/solve"
)

func tileSetFromString(grid gs.Grid, str string) gs.TileSet {
	lines := strings.Split(strings.Trim(str, "{}"), "|")
	var ts gs.TileSet
	for i, line := range lines {
		y := len(lines) - i - 1
		x := -1
		for {
			index := strings.IndexByte(line[x+1:], 'x')
			if index < 0 {
				break
			}
			x = x + 1 + index
			tileWithColor := *grid.TileAt(x, y)
			tileWithColor.Data.Color = 1
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

func testSolvePathsAbstract(t *testing.T, level string, x1, y1, x2, y2 int, solutions []string) {
	t.Helper()

	grid := solve.NewGridSolver(gs.MakeGridFromString(level, 2))

	var expectedSolutions []gs.TileSet
	for i := range solutions {
		expectedSolutions = append(expectedSolutions, tileSetFromString(grid.Grid, solutions[i]))
	}

	ch := grid.SolvePath(gs.TileCoord{X: x1, Y: y1}, gs.TileCoord{X: x2, Y: y2}, 1)
	var actualSolutions []gs.TileSet
	for ts := range ch {
		actualSolutions = append(actualSolutions, ts)
	}

	testUnorderedTilesetSliceEq(t, expectedSolutions, actualSolutions)
}

func TestSolvePaths_levelA1(t *testing.T) {
	const level = `[gA/ ] [    ] [    ] [g   ]`

	solutions := []string{"xxxx"}
	testSolvePathsAbstract(t, level, 0, 0, 3, 0, solutions)
}

func TestSolvePaths_levelA2(t *testing.T) {
	const level = `
[gA/ ] [----] [    ] [g   ]
[    ] [    ] [    ] [----]
`
	solutions := []string{"x xx|xxx "}
	testSolvePathsAbstract(t, level, 0, 1, 3, 1, solutions)
}

func TestSolvePaths_levelA3(t *testing.T) {
	const level = `
[----] [    ] [    ] [    ] [----]
[gA/ ] [    ] [  / ] [    ] [g   ]
[----] [    ] [    ] [    ] [----]
`
	solutions := []string{" xxx |xx xx|     ", "     |xx xx| xxx "}
	testSolvePathsAbstract(t, level, 0, 1, 4, 1, solutions)
}

func TestSolvePaths_levelA4(t *testing.T) {
	const level = `
[gA/ ] [    ] [g   ]
[    ] [    ] [ A/ ]
`
	solutions := []string{"x x|xxx"}
	testSolvePathsAbstract(t, level, 0, 1, 2, 1, solutions)
}

func TestSolvePaths_levelA5(t *testing.T) {
	const level = `
[    ] [ A/ ] [    ] [    ]
[    ] [    ] [ A/ ] [    ]
[gA/ ] [    ] [  / ] [g   ]
`
	solutions := []string{"xxx |x xx|x  x"}
	testSolvePathsAbstract(t, level, 0, 0, 3, 0, solutions)
}

func TestSolvePaths_levelA6(t *testing.T) {
	const level = `
[    ] [    ] [    ] [    ] [ A/ ] [    ] [    ] [    ]
[    ] [    ] [ A/ ] [    ] [ A/ ] [    ] [ A/ ] [    ]
[gA/ ] [    ] [    ] [    ] [ A/ ] [    ] [    ] [g   ]
`
	solutions := []string{
		"xxx xxx |x x x x |x xxx xx",
		"    xxx |xxx x x |x xxx xx",
		"xxx xxx |x x x xx|x xxx  x",
		"    xxx |xxx x xx|x xxx  x",
	}
	testSolvePathsAbstract(t, level, 0, 0, 7, 0, solutions)
}

func TestSolvePaths_levelA9(t *testing.T) {
	const level = `
[    ] [ A/ ] [    ] [    ] [ A/ ] [    ] [    ]
[gA/ ] [    ] [    ] [ A/ ] [    ] [    ] [g   ]
[    ] [    ] [ A/ ] [    ] [ A/ ] [    ] [    ]
`
	solutions := []string{"   xxxx|x xx  x|xxx    ", "   xxx |x xx xx|xxx    "}
	testSolvePathsAbstract(t, level, 0, 1, 6, 1, solutions)
}

func TestSolvePaths_basicColorNonePath(t *testing.T) {
	const level = `
[ A/ ] [    ] [    ] [ A/ ]
[g / ] [    ] [    ] [g   ]
[ A/ ] [    ] [    ] [ A/ ]
`
	grid := gs.MakeGridFromString(level, 2)
	gridSolver := solve.NewGridSolver(grid)
	solutionsChan := gridSolver.SolvePath(gs.TileCoord{X: 0, Y: 1}, gs.TileCoord{X: 3, Y: 1}, gs.ColorNone)
	var solutions []gs.TileSet
	for solution := range solutionsChan {
		solutions = append(solutions, solution)
	}
	if len(solutions) != 1 {
		t.Fatalf("solutions length expected to be 1 but was %d", len(solutions))
	}
	expected := gs.NewTileSet(
		gs.Tile{Coord: gs.TileCoord{X: 1, Y: 2}, Data: gs.TileData{Color: 1, Type: gs.TypeBlank}},
		gs.Tile{Coord: gs.TileCoord{X: 2, Y: 2}, Data: gs.TileData{Color: 1, Type: gs.TypeBlank}},

		gs.Tile{Coord: gs.TileCoord{X: 0, Y: 1}, Data: gs.TileData{Sticky: true, Type: gs.TypeGoal}},
		gs.Tile{Coord: gs.TileCoord{X: 1, Y: 1}, Data: gs.TileData{Type: gs.TypeBlank}},
		gs.Tile{Coord: gs.TileCoord{X: 2, Y: 1}, Data: gs.TileData{Type: gs.TypeBlank}},
		gs.Tile{Coord: gs.TileCoord{X: 3, Y: 1}, Data: gs.TileData{Type: gs.TypeGoal}},

		gs.Tile{Coord: gs.TileCoord{X: 1, Y: 0}, Data: gs.TileData{Color: 1, Type: gs.TypeBlank}},
		gs.Tile{Coord: gs.TileCoord{X: 2, Y: 0}, Data: gs.TileData{Color: 1, Type: gs.TypeBlank}},
	)

	if !expected.Eq(solutions[0]) {
		expectedGrid := grid.Clone()
		actualGrid := grid.Clone()
		expectedGrid.ApplyTileSet(expected)
		actualGrid.ApplyTileSet(solutions[0])
		t.Errorf("solutions not equal")
		t.Errorf("expected:\n%v", expectedGrid)
		t.Errorf("actual:\n%v", actualGrid)
	}
}
