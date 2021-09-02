package gridspech_test

import (
	"testing"

	gs "github.com/deanveloper/gridspech-go"
)

// has every TileType, and all are valid
const validTestGrid = `
2m1  2k   0    0    0  
2    2m2  0    0    0  
0    0    0    2j1  0  
0    0    0    2k   0  
0    0    0    0    0  
1    1/   1/   1    0  
1    _    0    1    0  
1m2  1e   0    1e   0  
`

// has every TileType, and all are invalid
const invalidTestGrid = `
0k    2     0     0     1   
2     0     0     1k    0   
0     0     0     1     0j1 
0m2   1     0     0     1   
0m1   0     0     0     0   
0m2   0     0     0     0   
0     _     0     1e    0   
2j1   2/j1  2j1   1e    1   
`

// MakeValidGrid returns a grid which contains a tile of every single Type and Color,
// and all tiles are valid.
func MakeValidGrid() gs.Grid {
	return gs.MakeGridFromStringNew(validTestGrid, 2)
}

// MakeInvalidGrid returns a grid which contains a tile of every single Type and Color,
// and all non-blank and non-hole tiles are invalid.
func MakeInvalidGrid() gs.Grid {
	return gs.MakeGridFromStringNew(invalidTestGrid, 2)
}

// TestMakeGridFromString tests creating a grid from a string.
func TestMakeGridFromString(t *testing.T) {
	tiles := MakeValidGrid().Tiles

	cases := []struct {
		Actual, Expected gs.TileData
	}{
		{tiles[0][0].Data, gs.TileData{Type: gs.TypeDot1, Color: 2}},
		{tiles[0][1].Data, gs.TileData{Type: gs.TypeBlank, Color: 2}},
		{tiles[1][0].Data, gs.TileData{Type: gs.TypeCrown, Color: 2}},
		{tiles[1][1].Data, gs.TileData{Type: gs.TypeDot2, Color: 2}},
		{tiles[1][7].Data, gs.TileData{Type: gs.TypeGoal, Color: 1}},
		{tiles[1][6].Data, gs.TileData{Type: gs.TypeHole}},
		{tiles[1][5].Data, gs.TileData{Type: gs.TypeBlank, Color: 1, Sticky: true}},
		{tiles[3][2].Data, gs.TileData{Type: gs.TypePlus, Color: 2}},
	}

	for _, testCase := range cases {
		if testCase.Expected != testCase.Actual {
			t.Errorf("\nexpected: %#v\ngot:      %#v\n", testCase.Expected, testCase.Actual)
		}
	}
}
