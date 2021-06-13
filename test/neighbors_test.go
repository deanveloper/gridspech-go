package gridspech_test

import (
	"testing"

	gs "github.com/deanveloper/gridspech-go"
)

// TestDirections tests the NorthOf, WestOf, EastOf, and SouthOf functions.
func TestDirections(t *testing.T) {
	grid := MakeValidGrid()
	tiles := grid.Tiles
	cases := []struct {
		Actual, Expected gs.Tile
	}{
		// test edges
		{grid.NorthOf(tiles[0][7]), gs.Tile{}},
		{grid.WestOf(tiles[0][7]), gs.Tile{}},
		{grid.EastOf(tiles[4][0]), gs.Tile{}},
		{grid.SouthOf(tiles[4][0]), gs.Tile{}},
		// test next to hole
		{grid.SouthOf(tiles[1][7]), gs.Tile{}},

		// test general cases
		{grid.NorthOf(tiles[1][1]), gs.Tile{Type: gs.TypeBlank}},
		{grid.WestOf(tiles[1][1]), gs.Tile{Type: gs.TypeBlank, Color: 2}},
		{grid.EastOf(tiles[1][1]), gs.Tile{Type: gs.TypeBlank}},
		{grid.SouthOf(tiles[1][1]), gs.Tile{Type: gs.TypeCrown, Color: 2}},
	}

	for _, testCase := range cases {
		// discard X/Y so we can use ==
		testCase.Actual.X = 0
		testCase.Actual.Y = 0

		if testCase.Expected != testCase.Actual {
			t.Errorf("\nexpected: %#v\ngot:      %#v\n", testCase.Expected, testCase.Actual)
		}
	}
}

// TestNeighbors tests the Neighbors function.
func TestNeighbors(t *testing.T) {
	grid := MakeValidGrid()
	tiles := grid.Tiles

	cases := []struct {
		Tile     gs.Tile
		Expected gs.TileSet
	}{
		// test corners
		{tiles[0][0], gs.NewTileSet(tiles[0][1], tiles[1][0])},
		{tiles[4][0], gs.NewTileSet(tiles[4][1], tiles[3][0])},
		{tiles[0][7], gs.NewTileSet(tiles[0][6], tiles[1][7])},
		{tiles[4][7], gs.NewTileSet(tiles[4][6], tiles[3][7])},

		// test next to hole
		{tiles[2][6], gs.NewTileSet(tiles[2][7], tiles[2][5], tiles[3][6])},

		// test tile with 4 neighbors
		{tiles[1][1], gs.NewTileSet(tiles[0][1], tiles[1][0], tiles[2][1], tiles[1][2])},
	}

	for _, testCase := range cases {
		actual := grid.Neighbors(testCase.Tile)
		if !actual.Eq(testCase.Expected) {
			t.Errorf("\nexpected: %#v\ngot:      %#v\n", testCase.Expected, actual)
		}
	}
}

func TestNeighborsWith(t *testing.T) {

	grid := MakeValidGrid()
	tiles := grid.Tiles

	noColor := func(t gs.Tile) bool { return t.Color == 0 }
	goalsOnly := func(t gs.Tile) bool { return t.Type == gs.TypeGoal }

	cases := []struct {
		Tile     gs.Tile
		Func     func(t gs.Tile) bool
		Expected gs.TileSet
	}{
		{tiles[0][0], noColor, gs.NewTileSet()},
		{tiles[1][4], noColor, gs.NewTileSet(tiles[0][4], tiles[1][3], tiles[2][4])},
		{tiles[0][0], goalsOnly, gs.NewTileSet()},
		{tiles[2][7], goalsOnly, gs.NewTileSet(tiles[1][7], tiles[3][7])},
	}

	for _, testCase := range cases {
		actual := grid.NeighborsWith(testCase.Tile, testCase.Func)
		if !actual.Eq(testCase.Expected) {
			t.Errorf("\nexpected: %#v\ngot:      %#v\n", testCase.Expected, actual)
		}
	}
}
