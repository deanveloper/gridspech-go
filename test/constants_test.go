package gridspech_test

import (
	"testing"

	gs "github.com/deanveloper/gridspech-go"
)

// has every TileType, and all are valid
const validTestGrid = `
[2A  ] [gA  ] [    ] [gA  ] [    ]
[ A  ] [----] [    ] [ A  ] [    ]
[ A  ] [ A/ ] [ A/ ] [ A  ] [    ]
[    ] [    ] [    ] [    ] [    ]
[    ] [    ] [    ] [cB  ] [    ]
[    ] [    ] [    ] [+B  ] [    ]
[ B  ] [2B  ] [    ] [    ] [    ]
[1B  ] [cB  ] [    ] [    ] [    ]
`

// has every TileType, and all are invalid
const invalidTestGrid = `
[+B  ] [+B/ ] [+B  ] [gA  ] [ A  ]
[    ] [----] [    ] [gA  ] [    ]
[2   ] [    ] [    ] [    ] [    ]
[1   ] [    ] [    ] [    ] [    ]
[2   ] [ A  ] [    ] [    ] [ A  ]
[    ] [    ] [    ] [ A  ] [+   ]
[ B  ] [    ] [    ] [cA  ] [    ]
[c   ] [ B  ] [    ] [    ] [ A  ]
`

// MakeValidGrid returns a grid which contains a tile of every single Type and Color,
// and all tiles are valid.
func MakeValidGrid() gs.Grid {
	return gs.MakeGridFromString(validTestGrid)
}

// MakeInvalidGrid returns a grid which contains a tile of every single Type and Color,
// and all non-blank and non-hole tiles are invalid.
func MakeInvalidGrid() gs.Grid {
	return gs.MakeGridFromString(invalidTestGrid)
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
