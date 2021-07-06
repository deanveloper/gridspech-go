package gridspech

import (
	"fmt"
	"strings"
)

// Width of the grid.
func (g Grid) Width() int {
	return len(g.Tiles)
}

// Height of the grid.
func (g Grid) Height() int {
	return len(g.Tiles[0])
}

// NorthOf returns the tile north of t in g.
func (g Grid) NorthOf(t Tile) Tile {
	if t.Y == g.Height()-1 || t.Type == TypeHole {
		return Tile{}
	}
	return g.Tiles[t.X][t.Y+1]
}

// EastOf returns the tile east of t in g.
func (g Grid) EastOf(t Tile) Tile {
	if t.X == g.Width()-1 || t.Type == TypeHole {
		return Tile{}
	}
	return g.Tiles[t.X+1][t.Y]
}

// SouthOf returns the tile south of t in g.
func (g Grid) SouthOf(t Tile) Tile {
	if t.Y == 0 || t.Type == TypeHole {
		return Tile{}
	}
	return g.Tiles[t.X][t.Y-1]
}

// WestOf returns the tile west of t in g.
func (g Grid) WestOf(t Tile) Tile {
	if t.X == 0 || t.Type == TypeHole {
		return Tile{}
	}
	return g.Tiles[t.X-1][t.Y]
}

// Neighbors returns all tiles directly next to t.
func (g Grid) Neighbors(t Tile) TileSet {
	var ts TileSet
	if neighbor := g.NorthOf(t); neighbor.Type != TypeHole {
		ts.Add(neighbor)
	}
	if neighbor := g.EastOf(t); neighbor.Type != TypeHole {
		ts.Add(neighbor)
	}
	if neighbor := g.SouthOf(t); neighbor.Type != TypeHole {
		ts.Add(neighbor)
	}
	if neighbor := g.WestOf(t); neighbor.Type != TypeHole {
		ts.Add(neighbor)
	}
	return ts
}

// TilesWith returns all non-hole tiles such that `pred` returns true.
func (g Grid) TilesWith(pred func(Tile) bool) TileSet {
	var ts TileSet

	for _, col := range g.Tiles {
		for _, tile := range col {
			if tile.Type != TypeHole && pred(tile) {
				ts.Add(tile)
			}
		}
	}

	return ts
}

// Blob returns all tiles which can form a path to t such that all tiles in the path have the same Color.
func (g Grid) Blob(t Tile) TileSet {
	var ts TileSet

	g.blobRecur(t, &ts)

	return ts
}

func (g Grid) blobRecur(t Tile, ts *TileSet) {
	neighbors := g.NeighborsWith(t, func(other Tile) bool {
		return other.Color == t.Color
	})

	for _, neighbor := range neighbors.Slice() {
		if !ts.Has(neighbor) {
			ts.Add(neighbor)
			g.blobRecur(neighbor, ts)
		}
	}
}

// NeighborsWith returns the set of neighbors such that `pred` returns true
func (g Grid) NeighborsWith(t Tile, pred func(o Tile) bool) TileSet {
	neighbors := g.Neighbors(t)
	for _, neighbor := range neighbors.Slice() {
		if !pred(neighbor) {
			neighbors.Remove(neighbor)
		}
	}
	return neighbors
}

// SetState sets the state of t in g.
func (g Grid) SetState(t Tile, state TileColor) {
	if !t.Sticky {
		g.Tiles[t.X][t.Y].Color = state
	}
}

// MakeGridFromString returns a Grid made from a string.
// See Grid.String() and Tile.String() for the format.
//
// May panic if the format is invalid.
func MakeGridFromString(str string) Grid {
	var grid Grid

	lines := strings.Split(strings.Trim(str, "\n"), "\n")

	height := len(lines)
	width := strings.Count(lines[0], "[")

	grid.Tiles = make([][]Tile, width)

	for x := 0; x < width; x++ {
		grid.Tiles[x] = make([]Tile, height)

		for y := 0; y < height; y++ {
			index := x * 6
			substr := lines[height-y-1][index+1 : index+4]
			typeByte, colorByte, stickyByte := substr[0], substr[1], substr[2]

			tile := grid.tileFromBytes(typeByte, colorByte, stickyByte)
			tile.X = x
			tile.Y = y
			grid.Tiles[x][y] = tile
		}
	}

	return grid
}

func (g Grid) tileFromBytes(typ, color, sticky byte) Tile {
	var tile Tile
	switch typ {
	case ' ':
		tile.Type = TypeBlank
	case 'g':
		tile.Type = TypeGoal
	case 'c':
		tile.Type = TypeCrown
	case '1':
		tile.Type = TypeDot1
	case '2':
		tile.Type = TypeDot2
	case '+':
		tile.Type = TypePlus
	}

	if color >= 'A' && color <= 'Z' {
		tile.Color = TileColor(color - 'A' + 1)
	}
	switch color {
	case 'O', ' ':
		tile.Color = 0
	case 'A':
		tile.Color = 1
	case 'B':
		tile.Color = 2
	}

	tile.Sticky = sticky == '/'

	return tile
}

func (t Tile) String() string {
	if t.Type == TypeHole {
		return "[---]"
	}

	var typeChar rune
	switch t.Type {
	case TypeBlank:
		typeChar = ' '
	case TypeGoal:
		typeChar = 'g'
	case TypeCrown:
		typeChar = 'c'
	case TypeDot1:
		typeChar = '1'
	case TypeDot2:
		typeChar = '2'
	}

	var colorChar rune
	switch t.Color {
	case ColorNone:
		colorChar = ' '
	case ColorA:
		colorChar = 'A'
	case ColorB:
		colorChar = 'B'
	}

	stickyChar := ' '
	if t.Sticky {
		stickyChar = '/'
	}
	return fmt.Sprintf("[%c%c%c]", typeChar, colorChar, stickyChar)
}

// String returns the string representation of g.
func (g Grid) String() string {
	byteSlice := make([]byte, (g.Width()*6)*g.Height()-1)
	for x, col := range g.Tiles {
		for y, tile := range col {
			index := x*6 + (g.Height()-y-1)*g.Width()*6

			copy(byteSlice[index:index+5], tile.String())
			if x < g.Width()-1 {
				byteSlice[index+5] = ' '
			}
			if x == g.Width()-1 && y != 0 {
				byteSlice[index+5] = '\n'
			}
		}
	}
	return string(byteSlice)
}

// Clone returns a clone of the grid. Modifications to the new grid will not modify the original grid.
func (g Grid) Clone() Grid {
	var newGrid Grid
	newGrid.Tiles = make([][]Tile, 0, len(g.Tiles))

	for _, col := range g.Tiles {
		newCol := make([]Tile, len(col))
		newGrid.Tiles = append(newGrid.Tiles, newCol)
		copy(newCol, col)
	}
	return newGrid
}
