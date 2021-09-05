package gridspech

import (
	"strings"
)

// TileSet represents a mathematical set of tiles. Tiles are compared using ==.
type TileSet struct {
	set map[Tile]struct{}
}

// NewTileSet returns a TileSet containing only tiles.
func NewTileSet(tiles ...Tile) TileSet {
	var ts TileSet
	for _, tile := range tiles {
		ts.Add(tile)
	}
	return ts
}

// Init initializes the tileset.
func (ts *TileSet) checkInit() {
	if ts.set == nil {
		ts.set = make(map[Tile]struct{})
	}
}

// Add adds t to the TileSet ts.
func (ts *TileSet) Add(t Tile) {
	ts.checkInit()
	ts.set[t] = struct{}{}
}

// Has returns if ts contains t.
func (ts TileSet) Has(t Tile) bool {
	_, ok := ts.set[t]
	return ok
}

// Remove removes t from ts.
func (ts *TileSet) Remove(t Tile) {
	ts.checkInit()
	delete(ts.set, t)
}

// RemoveIf removes each value for which pred returns true.
func (ts *TileSet) RemoveIf(pred func(t Tile) bool) {
	for tile := range ts.set {
		if pred(tile) {
			ts.Remove(tile)
		}
	}
}

// RemoveAll removes all of the elements in o from ts (making ts the intersection of ts and o)
func (ts *TileSet) RemoveAll(o TileSet) {
	if ts.Len() < o.Len() {
		for tile := range ts.set {
			if o.Has(tile) {
				ts.Remove(tile)
			}
		}
	} else {
		for tile := range o.set {
			if ts.Has(tile) {
				ts.Remove(tile)
			}
		}
	}
}

// Len returns the number of tiles in ts.
func (ts TileSet) Len() int {
	return len(ts.set)
}

// Merge adds all tiles in other into ts.
func (ts *TileSet) Merge(other TileSet) {
	ts.checkInit()
	for tile := range other.set {
		ts.set[tile] = struct{}{}
	}
}

// Eq returns if ts contains exactly the same contents as other.
func (ts TileSet) Eq(other TileSet) bool {
	if ts.Len() != other.Len() {
		return false
	}
	for tile := range ts.set {
		if !other.Has(tile) {
			return false
		}
	}
	return true
}

// Iter returns an iterator for this TileSet.
func (ts TileSet) Iter() <-chan Tile {
	iter := make(chan Tile, 5)

	go func() {
		for tile := range ts.set {
			iter <- tile
		}
		close(iter)
	}()

	return iter
}

// Slice returns a slice representation of ts
func (ts TileSet) Slice() []Tile {
	slice := make([]Tile, 0, len(ts.set))
	for tile := range ts.set {
		slice = append(slice, tile)
	}
	return slice
}

// ToTileCoordSet converts ts into a TileCoordSet
func (ts TileSet) ToTileCoordSet() TileCoordSet {
	var result TileCoordSet
	for val := range ts.set {
		result.Add(val.Coord)
	}
	return result
}

func (ts TileSet) String() string {
	slice := ts.Slice()

	var maxX, maxY int
	for _, tile := range slice {
		if tile.Coord.X > maxX {
			maxX = tile.Coord.X
		}
		if tile.Coord.Y > maxY {
			maxY = tile.Coord.Y
		}
	}
	maxX++
	maxY++

	tilesAt := make([][]Tile, maxX)
	for x := range tilesAt {
		tilesAt[x] = make([]Tile, maxY)
	}
	for _, v := range slice {
		tilesAt[v.Coord.X][v.Coord.Y] = v
	}

	var sb strings.Builder
	sb.WriteByte('{')
	for y := maxY - 1; y >= 0; y-- {
		for x := 0; x < maxX; x++ {
			if tile := tilesAt[x][y]; tile.Data.Type != TypeHole {
				sb.WriteByte(byte(tile.Data.Color) + '0')
			} else {
				sb.WriteByte(' ')
			}
		}
		if y > 0 {
			sb.WriteByte('|')
		}
	}
	sb.WriteByte('}')
	return sb.String()
}

// MultiLineString returns a string representation of this tileset on multiple lines
func (ts TileSet) MultiLineString() string {
	next := ts.String()
	next = next[1 : len(next)-1]
	next = strings.ReplaceAll(next, "|", "\n")
	next += "\n"
	return next
}
