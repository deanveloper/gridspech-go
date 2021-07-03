package solve_test

import (
	"strings"
	"testing"

	gs "github.com/deanveloper/gridspech-go"
	"github.com/deanveloper/gridspech-go/solve"
)

func tileSetFromString(grid gs.Grid, str string) gs.TileSet {
	lines := strings.Split(str, "\n")
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
			ts.Add(grid.Tiles[x][y])
		}
	}
	return ts
}

func tileSetHas(t *testing.T, ts gs.TileSet, tile gs.Tile) {
	t.Helper()

	if !ts.Has(tile) {
		t.Errorf(`expected tileset to have tile %v`, tile)
	}
}
func tileSetNotHas(t *testing.T, ts gs.TileSet, tile gs.Tile) {
	t.Helper()

	if ts.Has(tile) {
		t.Errorf(`expected tileset to not have tile %v`, tile)
	}
}

func testUnorderedTilesetSliceEq(t *testing.T, expected, actual []gs.TileSet) {
	t.Helper()

	if len(expected) != len(actual) {
		t.Errorf("not correct length, expected %d, actual %d\nexpected: %v\nactual: %v", len(expected), len(actual), expected, actual)
		return
	}

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
			if expected[i1].Eq(actual[i2]) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("incorrect solution %v", expected[i1])
		}
	}
}

func testSolvePathsAbstract(t *testing.T, level string, x1, y1, x2, y2 int, solutions []string) {
	t.Helper()

	grid := solve.Grid{Grid: gs.MakeGridFromString(level)}

	var expectedSolutions []gs.TileSet
	for i := range solutions {
		expectedSolutions = append(expectedSolutions, tileSetFromString(grid.Grid, solutions[i]))
	}

	ch := grid.SolveGoals(grid.Tiles[x1][y1], grid.Tiles[x2][y2])
	var actualSolutions []gs.TileSet
	for ts := range ch {
		actualSolutions = append(actualSolutions, ts)
	}
}

func TestSolvePaths_levelA1(t *testing.T) {
	const level = `[gA/] [   ] [   ] [g  ]`

	solutions := []string{"xxxx"}
	testSolvePathsAbstract(t, level, 0, 0, 3, 0, solutions)
}

func TestSolvePaths_level2(t *testing.T) {
	const level = `
[gA/]       [   ] [g  ]
[   ] [   ] [   ]      
`
	solutions := []string{"x xx\nxxx "}
	testSolvePathsAbstract(t, level, 0, 1, 3, 1, solutions)
}

func TestSolvePaths_level3(t *testing.T) {
	const level = `
      [   ] [   ] [   ]      
[gA/] [   ] [  /] [   ] [g  ]
      [   ] [   ] [   ]      
`
	solutions := []string{" xxx \nxx xx\n     ", "     \nxx xx\n xxx "}
	testSolvePathsAbstract(t, level, 0, 1, 4, 1, solutions)
}

func TestSolvePaths_level4(t *testing.T) {
	const level = `
[gA/] [   ] [g  ]
[   ] [   ] [ A/]
`
	solutions := []string{"x x\nxxx"}
	testSolvePathsAbstract(t, level, 0, 1, 2, 1, solutions)
}

func TestSolvePaths_level5(t *testing.T) {
	const level = `
[   ] [ A/] [   ] [   ]
[   ] [   ] [ A/] [   ]
[gA/] [   ] [  /] [g  ]
`
	solutions := []string{"xxx \nx xx\nx  x"}
	testSolvePathsAbstract(t, level, 0, 0, 3, 0, solutions)
}
