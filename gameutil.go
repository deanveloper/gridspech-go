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

// TileAt returns a reference to the tile at the coordinate.
func (g Grid) TileAt(x, y int) *Tile {
	return &g.Tiles[x][y]
}

// TileAtCoord returns a reference to the tile at the coordinate.
func (g Grid) TileAtCoord(coord TileCoord) *Tile {
	return &g.Tiles[coord.X][coord.Y]
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
func (g Grid) Neighbors(coord TileCoord) TileSet {
	t := *g.TileAtCoord(coord)
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

// NeighborsWith returns the set of neighbors such that `pred` returns true
func (g Grid) NeighborsWith(coord TileCoord, pred func(o Tile) bool) TileSet {
	neighbors := g.Neighbors(coord)
	for _, neighbor := range neighbors.Slice() {
		if !pred(neighbor) {
			neighbors.Remove(neighbor)
		}
	}
	return neighbors
}

// TilesWith returns all non-hole tiles such that `pred` returns true.
func (g Grid) TilesWith(pred func(o Tile) bool) TileSet {
	var ts TileSet

	for _, col := range g.Tiles {
		for _, tile := range col {
			if tile.Data.Type != TypeHole && pred(tile) {
				ts.Add(tile)
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
	t := g.TileAtCoord(coord)
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

func (t Tile) String() string {
	return t.Data.String()
}

func (td TileData) String() string {
	if td.Type == TypeHole {
		return "_"
	}

	// format:
	// `${t.value}${lock}${sym}${wrapl}${wrapu}${wrapd}${wrapr}`
	var sb strings.Builder

	sb.WriteByte(byte(td.Color) + '0')
	if td.Sticky {
		sb.WriteByte('/')
	}
	switch td.Type {
	case TypeBlank:
		break
	case TypeCrown:
		sb.WriteByte('k')
	case TypeGoal:
		sb.WriteByte('e')
	case TypeDot1:
		sb.WriteString("m1")
	case TypeDot2:
		sb.WriteString("m2")
	case TypeDot3:
		sb.WriteString("m3")
	case TypePlus:
		sb.WriteString("j1")
	default:
		panic(fmt.Sprintf("invalid type %d", td.Type))
	}

	if td.ArrowWest {
		sb.WriteByte('<')
	}
	if td.ArrowNorth {
		sb.WriteByte('^')
	}
	if td.ArrowSouth {
		sb.WriteByte('v')
	}
	if td.ArrowEast {
		sb.WriteByte('>')
	}
	return sb.String()
}

// Clone returns a clone of the grid. Modifications to the new grid will not modify the original grid.
func (g Grid) Clone() Grid {
	var newGrid Grid
	newGrid.Tiles = make([][]Tile, 0, len(g.Tiles))
	newGrid.MaxColors = g.MaxColors

	for _, col := range g.Tiles {
		newCol := make([]Tile, len(col))
		newGrid.Tiles = append(newGrid.Tiles, newCol)
		copy(newCol, col)
	}
	return newGrid
}
