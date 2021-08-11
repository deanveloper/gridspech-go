package gridspech

import (
	"fmt"
	"strings"
)

var wordsToArrowsMap = map[string]rune{
	"north":              '╵',
	"northeast":          '└',
	"northsouth":         '│',
	"northwest":          '┘',
	"northeastsouth":     '├',
	"northeastwest":      '┴',
	"northsouthwest":     '┤',
	"northeastsouthwest": '┼',
	"east":               '╶',
	"eastsouth":          '┌',
	"eastwest":           '─',
	"eastsouthwest":      '┬',
	"south":              '╷',
	"southwest":          '┐',
	"west":               '╸',
}
var arrowsToWordsMap = map[rune]string{
	'╵': "north",
	'└': "northeast",
	'│': "northsouth",
	'┘': "northwest",
	'├': "northeastsouth",
	'┴': "northeastwest",
	'┤': "northsouthwest",
	'┼': "northeastsouthwest",
	'╶': "east",
	'┌': "eastsouth",
	'─': "eastwest",
	'┬': "eastsouthwest",
	'╷': "south",
	'┐': "southwest",
	'╸': "west",
}

func encodeArrows(north, east, south, west bool) rune {
	var str strings.Builder
	if north {
		str.WriteString("north")
	}
	if east {
		str.WriteString("east")
	}
	if south {
		str.WriteString("south")
	}
	if west {
		str.WriteString("west")
	}
	if str.Len() == 0 {
		return ' '
	}
	return wordsToArrowsMap[str.String()]
}
func decodeArrows(arrows rune) (north, east, south, west bool) {
	oldStr := arrowsToWordsMap[arrows]
	nextStr := oldStr
	if nextStr = strings.TrimPrefix(oldStr, "north"); oldStr != nextStr {
		north = true
	}
	oldStr = nextStr
	if nextStr = strings.TrimPrefix(oldStr, "east"); oldStr != nextStr {
		east = true
	}
	oldStr = nextStr
	if nextStr = strings.TrimPrefix(oldStr, "south"); oldStr != nextStr {
		south = true
	}
	oldStr = nextStr
	if nextStr = strings.TrimPrefix(oldStr, "west"); oldStr != nextStr {
		west = true
	}

	return north, east, south, west
}

// Width of the grid.
func (g Grid) Width() int {
	return len(g.Tiles)
}

// Height of the grid.
func (g Grid) Height() int {
	return len(g.Tiles[0])
}

// TileAt returns the tile at a specific coordinate.
func (g Grid) TileAt(coord TileCoord) Tile {
	return g.Tiles[coord.X][coord.Y]
}

// NorthOf returns the tile north of t in g.
func (g Grid) NorthOf(t Tile) Tile {
	// special behavior if we have an arrow
	if t.Data.ArrowNorth {
		nextTile := findTile(t, func(each Tile) (Tile, bool) {
			if each.Coord.Y == g.Height()-1 {
				return g.Tiles[each.Coord.X][0], true
			}
			if each.Data.Type == TypeHole {
				return g.Tiles[each.Coord.X][each.Coord.Y+1], true
			}
			return Tile{}, false
		})
		return nextTile
	}

	if t.Coord.Y == g.Height()-1 || t.Data.Type == TypeHole {
		return Tile{}
	}
	return g.Tiles[t.Coord.X][t.Coord.Y+1]
}

// EastOf returns the tile east of t in g.
func (g Grid) EastOf(t Tile) Tile {
	// special behavior if we have an arrow
	if t.Data.ArrowEast {
		nextTile := findTile(t, func(each Tile) (Tile, bool) {
			if each.Coord.X == g.Width()-1 {
				return g.Tiles[0][each.Coord.Y], true
			}
			if each.Data.Type == TypeHole {
				return g.Tiles[each.Coord.X+1][each.Coord.Y], true
			}
			return Tile{}, false
		})
		return nextTile
	}

	if t.Coord.X == g.Width()-1 || t.Data.Type == TypeHole {
		return Tile{}
	}
	return g.Tiles[t.Coord.X+1][t.Coord.Y]
}

// SouthOf returns the tile south of t in g.
func (g Grid) SouthOf(t Tile) Tile {
	// special behavior if we have an arrow
	if t.Data.ArrowSouth {
		nextTile := findTile(t, func(each Tile) (Tile, bool) {
			if each.Coord.Y == 0 {
				return g.Tiles[each.Coord.X][g.Height()-1], true
			}
			if each.Data.Type == TypeHole {
				return g.Tiles[each.Coord.X][each.Coord.Y-1], true
			}
			return Tile{}, false
		})
		return nextTile
	}

	if t.Coord.Y == 0 || t.Data.Type == TypeHole {
		return Tile{}
	}
	return g.Tiles[t.Coord.X][t.Coord.Y-1]
}

// WestOf returns the tile west of t in g.
func (g Grid) WestOf(t Tile) Tile {
	// special behavior if we have an arrow
	if t.Data.ArrowWest {
		nextTile := findTile(t, func(each Tile) (Tile, bool) {
			if each.Coord.X == 0 {
				return g.Tiles[g.Width()-1][each.Coord.Y], true
			}
			if each.Data.Type == TypeHole {
				return g.Tiles[each.Coord.X-1][each.Coord.Y], true
			}
			return Tile{}, false
		})
		return nextTile
	}

	if t.Coord.X == 0 || t.Data.Type == TypeHole {
		return Tile{}
	}
	return g.Tiles[t.Coord.X-1][t.Coord.Y]
}

func findTile(start Tile, iter func(each Tile) (Tile, bool)) Tile {
	last := start
	next, ok := iter(start)
	for {
		if !ok {
			return last
		}
		if next == start {
			panic("no tile found")
		}
		last = next
		next, ok = iter(next)
	}
}

// Neighbors returns all tiles directly next to t.
func (g Grid) Neighbors(t Tile) TileSet {
	var ts TileSet
	if neighbor := g.NorthOf(t); neighbor.Data.Type != TypeHole {
		ts.Add(neighbor)
	}
	if neighbor := g.EastOf(t); neighbor.Data.Type != TypeHole {
		ts.Add(neighbor)
	}
	if neighbor := g.SouthOf(t); neighbor.Data.Type != TypeHole {
		ts.Add(neighbor)
	}
	if neighbor := g.WestOf(t); neighbor.Data.Type != TypeHole {
		ts.Add(neighbor)
	}
	return ts
}

// TilesWith returns all non-hole tiles such that `pred` returns true.
func (g Grid) TilesWith(pred func(o Tile) bool) TileCoordSet {
	var ts TileCoordSet

	for _, col := range g.Tiles {
		for _, tile := range col {
			if tile.Data.Type != TypeHole && pred(tile) {
				ts.Add(tile.Coord)
			}
		}
	}

	return ts
}

// ApplyTileSet will loop through ts and update all tiles
// with the same coordinates to have the same data as the tiles in ts.
func (g Grid) ApplyTileSet(ts TileSet) {
	for _, tile := range ts.Slice() {
		g.Tiles[tile.Coord.X][tile.Coord.Y] = tile
	}
}

// Blob returns all tiles which can form a path to t such that all tiles in the path have the same Color.
func (g Grid) Blob(coord TileCoord) TileSet {
	var ts TileSet

	g.blobRecur(coord, &ts)

	return ts
}

func (g Grid) blobRecur(coord TileCoord, ts *TileSet) {
	t := g.TileAt(coord)
	neighbors := g.NeighborsWith(coord, func(o Tile) bool {
		return o.Data.Color == t.Data.Color
	})

	for _, neighbor := range neighbors.Slice() {
		if !ts.Has(neighbor) {
			ts.Add(neighbor)
			g.blobRecur(neighbor.Coord, ts)
		}
	}
}

// NeighborsWith returns the set of neighbors such that `pred` returns true
func (g Grid) NeighborsWith(coord TileCoord, pred func(o Tile) bool) TileSet {
	t := g.TileAt(coord)
	neighbors := g.Neighbors(t)
	for _, neighbor := range neighbors.Slice() {
		if !pred(neighbor) {
			neighbors.Remove(neighbor)
		}
	}
	return neighbors
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
			index := x * 7
			substr := []rune(lines[height-y-1][index+1 : index+5])
			typeRune, colorRune, stickyRune, arrowsRune := substr[0], substr[1], substr[2], substr[3]

			tile := tileFromRunes(typeRune, colorRune, stickyRune, arrowsRune)
			tile.Coord.X = x
			tile.Coord.Y = y
			grid.Tiles[x][y] = tile
		}
	}

	return grid
}

func tileFromRunes(typ, color, sticky, arrows rune) Tile {
	var tile Tile
	switch typ {
	case ' ':
		tile.Data.Type = TypeBlank
	case 'g':
		tile.Data.Type = TypeGoal
	case 'c':
		tile.Data.Type = TypeCrown
	case '1':
		tile.Data.Type = TypeDot1
	case '2':
		tile.Data.Type = TypeDot2
	case '3':
		tile.Data.Type = TypeDot3
	case '+':
		tile.Data.Type = TypePlus
	}

	if color >= 'A' && color <= 'Z' {
		tile.Data.Color = TileColor(color - 'A' + 1)
	}
	switch color {
	case 'O', ' ':
		tile.Data.Color = 0
	case 'A':
		tile.Data.Color = 1
	case 'B':
		tile.Data.Color = 2
	}

	tile.Data.Sticky = sticky == '/'

	tile.Data.ArrowNorth, tile.Data.ArrowEast, tile.Data.ArrowSouth, tile.Data.ArrowWest = decodeArrows(arrows)

	return tile
}

func (t Tile) String() string {
	return t.Data.String()
}

func (td TileData) String() string {
	if td.Type == TypeHole {
		return "[----]"
	}

	var typeChar rune
	switch td.Type {
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
	case TypeDot3:
		typeChar = '3'
	default:
		panic(fmt.Sprint("invalid Type", td.Type))
	}

	var colorChar rune
	switch td.Color {
	case 0:
		colorChar = ' '
	case 1:
		colorChar = 'A'
	case 2:
		colorChar = 'B'

		// special case for "unknown" color in solvers
	case 100:
		colorChar = ' '
	default:
		panic(fmt.Sprintf("invalid color %d", td.Color))
	}

	stickyChar := ' '
	if td.Sticky {
		stickyChar = '/'
	}

	arrowsChar := encodeArrows(td.ArrowNorth, td.ArrowEast, td.ArrowSouth, td.ArrowWest)

	str := fmt.Sprintf("[%c%c%c%c]", typeChar, colorChar, stickyChar, arrowsChar)
	return str
}

// String returns the string representation of g.
func (g Grid) String() string {
	byteSlice := make([]byte, (g.Width()*7)*g.Height()-1)
	for x, col := range g.Tiles {
		for y, tile := range col {
			index := x*7 + (g.Height()-y-1)*g.Width()*7

			copy(byteSlice[index:index+6], tile.String())
			if x < g.Width()-1 {
				byteSlice[index+6] = ' '
			}
			if x == g.Width()-1 && y != 0 {
				byteSlice[index+6] = '\n'
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
