package gridspech

import "strings"

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
	if t.Y == g.Height()-1 || t.Type == Hole {
		return Tile{}
	}
	return g.Tiles[t.X][t.Y+1]
}

// EastOf returns the tile east of t in g.
func (g Grid) EastOf(t Tile) Tile {
	if t.X == g.Width()-1 || t.Type == Hole {
		return Tile{}
	}
	return g.Tiles[t.X+1][t.Y]
}

// SouthOf returns the tile south of t in g.
func (g Grid) SouthOf(t Tile) Tile {
	if t.Y == 0 || t.Type == Hole {
		return Tile{}
	}
	return g.Tiles[t.X][t.Y-1]
}

// WestOf returns the tile west of t in g.
func (g Grid) WestOf(t Tile) Tile {
	if t.X == 0 || t.Type == Hole {
		return Tile{}
	}
	return g.Tiles[t.X-1][t.Y]
}

// Neighbors returns all tiles directly next to t.
func (g Grid) Neighbors(t Tile) TileSet {
	var ts TileSet
	ts.Init()
	if neighbor := g.NorthOf(t); neighbor.Type != Hole {
		ts.Add(neighbor)
	}
	if neighbor := g.EastOf(t); neighbor.Type != Hole {
		ts.Add(neighbor)
	}
	if neighbor := g.SouthOf(t); neighbor.Type != Hole {
		ts.Add(neighbor)
	}
	if neighbor := g.WestOf(t); neighbor.Type != Hole {
		ts.Add(neighbor)
	}
	return ts
}

// TilesWith returns all non-hole tiles such that `pred` returns true.
func (g Grid) TilesWith(pred func(Tile) bool) TileSet {
	var ts TileSet
	ts.Init()

	for _, col := range g.Tiles {
		for _, tile := range col {
			if tile.Type != Hole && pred(tile) {
				ts.Add(tile)
			}
		}
	}

	return ts
}

// Blob returns all tiles which can form a path to t such that all tiles in the path have the same state.
func (g Grid) Blob(t Tile) TileSet {
	var ts TileSet
	ts.Init()

	g.blobRecur(t, ts)

	return ts
}

func (g Grid) blobRecur(t Tile, ts TileSet) {
	neighbors := g.NeighborsWithState(t, t.State)

	for _, neighbor := range neighbors.Slice() {
		if !ts.Has(neighbor) {
			ts.Add(neighbor)
			g.blobRecur(neighbor, ts)
		}
	}
}

// NeighborsWith returns the set of neighbors such that `pred` returns true
func (g Grid) NeighborsWith(t Tile, pred func(Tile) bool) TileSet {
	neighbors := g.Neighbors(t)
	for _, neighbor := range neighbors.Slice() {
		if !pred(neighbor) {
			neighbors.Remove(neighbor)
		}
	}
	return neighbors
}

// NeighborsWithState returns the set of neighbors which have a certain state
func (g Grid) NeighborsWithState(t Tile, state TileState) TileSet {
	neighbors := g.Neighbors(t)
	for _, neighbor := range neighbors.Slice() {
		if neighbor.State != state {
			neighbors.Remove(neighbor)
		}
	}
	return neighbors
}

// SetState sets the state of t in g.
func (g Grid) SetState(t Tile, state TileState) {
	if !t.Sticky {
		g.Tiles[t.X][t.Y].State = state
	}
}

// MakeGridFromString returns a Grid made from a string.
// See Grid.String() for the format.
func MakeGridFromString(str string) Grid {
	var grid Grid
	firstLine := strings.SplitN(str, "\n", 1)[0]

	height := strings.Count(str, "\n")
	width := strings.Count(firstLine, "[")

	grid.Tiles = make([][]Tile, width)

	for x := 0; x < width; x++ {
		grid.Tiles[x] = make([]Tile, height)

		for y := 0; y < height; y++ {
			index := x*6 + (height-y)*width*6
			substr := str[index : index+4]
			holeByte, typeByte, stateByte, stickyByte := substr[0], substr[1], substr[2], substr[3]

			tile := grid.tileFromBytes(holeByte, typeByte, stateByte, stickyByte)
			tile.X = x
			tile.Y = y
			grid.Tiles[x][y] = tile
		}
	}

	return grid
}

func (g Grid) tileFromBytes(hole, typ, state, sticky byte) Tile {
	if hole == ' ' {
		return Tile{}
	}

	var tile Tile
	switch typ {
	case ' ':
		tile.Type = Blank
	case 'g':
		tile.Type = Goal
	case 'c':
		tile.Type = Crown
	case '1':
		tile.Type = Dot1
	case '2':
		tile.Type = Dot2
	}

	switch state {
	case 'O':
		tile.State = Disabled
	case 'X':
		tile.State = Enabled
	}

	tile.Sticky = sticky == '/'

	return tile
}

// String returns the string representation of g.
func (g Grid) String() string {
	var sb strings.Builder
	for _, col := range g.Tiles {
		for _, tile := range col {

			if tile.Type == Hole {
				sb.WriteString("[   ] ")
				continue
			}

			sb.WriteByte('[')

			switch tile.Type {
			case Blank:
				sb.WriteByte(' ')
			case Goal:
				sb.WriteByte('g')
			case Crown:
				sb.WriteByte('c')
			case Dot1:
				sb.WriteByte('1')
			case Dot2:
				sb.WriteByte('2')
			}

			switch tile.State {
			case Enabled:
				sb.WriteByte('X')
			case Disabled:
				sb.WriteByte('O')
			}
			if tile.Sticky {
				sb.WriteByte('/')
			} else {
				sb.WriteByte(' ')
			}

			sb.WriteByte(']')
			sb.WriteByte(' ')
		}
		sb.WriteByte('\n')
	}
	return sb.String()
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
