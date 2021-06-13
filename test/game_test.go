package game_test

import (
	"testing"

	gs "github.com/deanveloper/gridspech-go"
)

// has every TileType, and all are valid
var validTestGrid = `
[2A ] [gA ] [   ] [gA ] [   ]
[ A ]       [   ] [ A ] [   ]
[ A ] [ A/] [ A/] [ A ] [   ]
[   ] [   ] [   ] [   ] [   ]
[   ] [   ] [   ] [cB ] [   ]
[   ] [   ] [   ] [+B ] [   ]
[ B ] [2B ] [   ] [   ] [   ]
[1B ] [cB ] [   ] [   ] [   ]
`

// has every TileType, and all are invalid
var invalidTestGrid = `
[+B ] [+B/] [+B ] [gA ] [ A ]
[   ]       [   ] [gA ] [   ]
[2  ] [   ] [   ] [   ] [   ]
[1  ] [   ] [   ] [   ] [   ]
[2  ] [ A ] [   ] [   ] [ A ]
[   ] [   ] [   ] [ A ] [+  ]
[ B ] [   ] [   ] [cA ] [   ]
[c  ] [ B ] [   ] [   ] [ A ]
`

func MakeValidGrid() gs.Grid {
	return gs.MakeGridFromString(validTestGrid)
}
func MakeInvalidGrid() gs.Grid {
	return gs.MakeGridFromString(invalidTestGrid)
}

func TestMakeGridFromString(t *testing.T) {
	tiles := MakeValidGrid().Tiles

	cases := []struct {
		Actual, Expected gs.Tile
	}{
		{tiles[0][0], gs.Tile{Type: gs.TypeDot1, Color: 2}},
		{tiles[0][1], gs.Tile{Type: gs.TypeBlank, Color: 2}},
		{tiles[1][0], gs.Tile{Type: gs.TypeCrown, Color: 2}},
		{tiles[1][1], gs.Tile{Type: gs.TypeDot2, Color: 2}},
		{tiles[1][7], gs.Tile{Type: gs.TypeGoal, Color: 1}},
		{tiles[1][6], gs.Tile{Type: gs.TypeHole}},
		{tiles[1][5], gs.Tile{Type: gs.TypeBlank, Color: 1, Sticky: true}},
		{tiles[3][2], gs.Tile{Type: gs.TypePlus, Color: 2}},
	}

	for _, testCase := range cases {
		testCase.Expected.X = testCase.Actual.X
		testCase.Expected.Y = testCase.Actual.Y

		if testCase.Expected != testCase.Actual {
			t.Errorf("\nexpected: %#v\ngot:      %#v\n", testCase.Expected, testCase.Actual)
		}
	}
}

func TestRules(t *testing.T) {

}
