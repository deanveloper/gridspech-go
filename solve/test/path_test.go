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

func TestSolvePaths_level1(t *testing.T) {
	const level = `[gA/] [   ] [   ] [g  ]`
	grid := solve.Grid{Grid: gs.MakeGridFromString(level)}
	solution := tileSetFromString(grid.Grid, "xxxx")

	// should return the one and only solution
	var sols int
	ch := grid.SolveGoals(grid.Tiles[0][0], grid.Tiles[3][0])
	for each := range ch {
		if !each.Eq(solution) {
			t.Errorf(`incorrect solution %v`, each)
		}
		sols++
	}

	if sols != 1 {
		t.Errorf(`expected %v solutions, found %v solutions`, 1, sols)
	}
}

func TestSolvePaths_level2(t *testing.T) {
	const level = `
[gA/]       [   ] [g  ]
[   ] [   ] [   ]      
`
	grid := solve.Grid{Grid: gs.MakeGridFromString(level)}
	solution := tileSetFromString(grid.Grid, "x xx\nxxx ")

	// should return the one and only solution
	var sols int
	ch := grid.SolveGoals(grid.Tiles[0][1], grid.Tiles[3][1])
	for each := range ch {
		if !each.Eq(solution) {
			t.Errorf(`incorrect solution %v`, each)
		}

		sols++
	}

	if sols != 1 {
		t.Errorf(`expected %v solutions, found %v solutions`, 1, sols)
	}
}
