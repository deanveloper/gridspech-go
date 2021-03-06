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
		y := t.Coord.Y
		for {
			y = (y + 1) % g.Height()
			if y == t.Coord.Y {
				break
			}

			if g.TileAt(t.Coord.X, y).Data.Type != TypeHole {
				break
			}
		}
		return *g.TileAt(t.Coord.X, y)
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
		x := t.Coord.X
		for {
			x = (x + 1) % g.Width()
			if x == t.Coord.X {
				break
			}

			if g.TileAt(x, t.Coord.Y).Data.Type != TypeHole {
				break
			}
		}
		return *g.TileAt(x, t.Coord.Y)
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
		y := t.Coord.Y
		for {
			y = (y + g.Height() - 1) % g.Height()
			if y == t.Coord.Y {
				break
			}

			if g.TileAt(t.Coord.X, y).Data.Type != TypeHole {
				break
			}
		}
		return *g.TileAt(t.Coord.X, y)
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
		x := t.Coord.X
		for {
			x = (x + g.Width() - 1) % g.Width()
			if x == t.Coord.X {
				break
			}

			if g.TileAt(x, t.Coord.Y).Data.Type != TypeHole {
				break
			}
		}
		return *g.TileAt(x, t.Coord.Y)
	}

	if t.Coord.X == 0 || t.Data.Type == TypeHole {
		return Tile{}
	}
	return g.Tiles[t.Coord.X-1][t.Coord.Y]
}

func findTile(start Tile, iter func(each Tile) (next Tile, ok bool)) Tile {
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

// NeighborSlice returns a slice of the tiles directly next to t.
// Note that arrows can cause a tile to be its own neighbor, or to
// have a tile appear more than once in the slice.
func (g Grid) NeighborSlice(coord TileCoord) []Tile {
	t := *g.TileAtCoord(coord)
	var neighbors []Tile
	if neighbor := g.NorthOf(t); neighbor.Data.Type != TypeHole {
		neighbors = append(neighbors, neighbor)
	}
	if neighbor := g.EastOf(t); neighbor.Data.Type != TypeHole {
		neighbors = append(neighbors, neighbor)
	}
	if neighbor := g.SouthOf(t); neighbor.Data.Type != TypeHole {
		neighbors = append(neighbors, neighbor)
	}
	if neighbor := g.WestOf(t); neighbor.Data.Type != TypeHole {
		neighbors = append(neighbors, neighbor)
	}
	return neighbors
}

// NeighborSliceWith returns a slice of the tiles directly next to t, such that pred returns true.
// Note that arrows can cause a tile to be its own neighbor, or to
// have a tile appear more than once in the slice.
func (g Grid) NeighborSliceWith(coord TileCoord, pred func(o Tile) bool) []Tile {
	t := *g.TileAtCoord(coord)
	var neighbors []Tile
	if neighbor := g.NorthOf(t); neighbor.Data.Type != TypeHole && pred(neighbor) {
		neighbors = append(neighbors, neighbor)
	}
	if neighbor := g.EastOf(t); neighbor.Data.Type != TypeHole && pred(neighbor) {
		neighbors = append(neighbors, neighbor)
	}
	if neighbor := g.SouthOf(t); neighbor.Data.Type != TypeHole && pred(neighbor) {
		neighbors = append(neighbors, neighbor)
	}
	if neighbor := g.WestOf(t); neighbor.Data.Type != TypeHole && pred(neighbor) {
		neighbors = append(neighbors, neighbor)
	}
	return neighbors
}

// NeighborSet returns a set of all tiles directly next to t.
// Note that because this is a TileSet, if a neighbor would appear
// more than once in NeighborSlice(), it will only appear once here.
func (g Grid) NeighborSet(coord TileCoord) TileSet {
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

// NeighborSetWith returns the set of neighbors such that `pred` returns true
func (g Grid) NeighborSetWith(coord TileCoord, pred func(o Tile) bool) TileSet {
	t := *g.TileAtCoord(coord)
	var ts TileSet
	if neighbor := g.NorthOf(t); neighbor.Data.Type != TypeHole && pred(neighbor) {
		ts.Add(neighbor)
	}
	if neighbor := g.EastOf(t); neighbor.Data.Type != TypeHole && pred(neighbor) {
		ts.Add(neighbor)
	}
	if neighbor := g.SouthOf(t); neighbor.Data.Type != TypeHole && pred(neighbor) {
		ts.Add(neighbor)
	}
	if neighbor := g.WestOf(t); neighbor.Data.Type != TypeHole && pred(neighbor) {
		ts.Add(neighbor)
	}
	return ts
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

	g.blobRecur(coord, &ts, func(o Tile) bool { return true })

	return ts
}

// BlobWith is similar to Blob, but with a filter function to exclude certain tiles.
func (g Grid) BlobWith(coord TileCoord, filter func(o Tile) bool) TileSet {
	var ts TileSet

	g.blobRecur(coord, &ts, filter)

	return ts
}

func (g Grid) blobRecur(coord TileCoord, ts *TileSet, filter func(o Tile) bool) {
	t := *g.TileAtCoord(coord)
	ts.Add(t)

	neighbors := g.NeighborSetWith(coord, func(o Tile) bool {
		return o.Data.Color == t.Data.Color && filter(o)
	})

	for _, neighbor := range neighbors.Slice() {
		if !ts.Has(neighbor) {
			g.blobRecur(neighbor.Coord, ts, filter)
		}
	}
}

func (t Tile) String() string {
	return fmt.Sprintf("[%v: %v]", t.Coord, t.Data)
}

func (t TileCoord) String() string {
	return fmt.Sprintf("(%d, %d)", t.X, t.Y)
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
	case TypeJoin1:
		sb.WriteString("j1")
	case TypeJoin2:
		sb.WriteString("j2")
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
