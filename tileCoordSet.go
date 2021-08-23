package gridspech

import "strings"

// TileCoordSet represents a mathematical set of coordinates.
type TileCoordSet struct {
	set map[TileCoord]struct{}
}

// NewTileCoordSet returns a TileCoordSet containing only tiles.
func NewTileCoordSet(tiles ...TileCoord) TileCoordSet {
	var cs TileCoordSet
	for _, tile := range tiles {
		cs.Add(tile)
	}
	return cs
}

// Init initializes the tilecoordset.
func (ts *TileCoordSet) checkInit() {
	if ts.set == nil {
		ts.set = make(map[TileCoord]struct{})
	}
}

// Add adds t to the TileSet ts.
func (ts *TileCoordSet) Add(t TileCoord) {
	ts.checkInit()
	ts.set[t] = struct{}{}
}

// Has returns if ts contains t.
func (ts TileCoordSet) Has(t TileCoord) bool {
	_, ok := ts.set[t]
	return ok
}

// Remove removes t from ts.
func (ts *TileCoordSet) Remove(t TileCoord) {
	ts.checkInit()
	delete(ts.set, t)
}

// RemoveIf removes each value for which pred returns true.
func (ts *TileCoordSet) RemoveIf(pred func(coord TileCoord) bool) {
	for tile := range ts.set {
		if pred(tile) {
			ts.Remove(tile)
		}
	}
}

// RemoveAll removes all of the elements in o from ts (making ts the intersection of ts and o)
func (ts *TileCoordSet) RemoveAll(o TileCoordSet) {
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
func (ts TileCoordSet) Len() int {
	return len(ts.set)
}

// Merge adds all tiles in other into ts.
func (ts *TileCoordSet) Merge(other TileCoordSet) {
	ts.checkInit()
	for tile := range other.set {
		ts.set[tile] = struct{}{}
	}
}

// Eq returns if ts contains exactly the same contents as other.
func (ts TileCoordSet) Eq(other TileCoordSet) bool {
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
func (ts TileCoordSet) Iter() <-chan TileCoord {
	iter := make(chan TileCoord, 5)

	go func() {
		for tile := range ts.set {
			iter <- tile
		}
		close(iter)
	}()

	return iter
}

// Slice returns a slice representation of ts
func (ts TileCoordSet) Slice() []TileCoord {
	slice := make([]TileCoord, 0, len(ts.set))
	for tile := range ts.set {
		slice = append(slice, tile)
	}
	return slice
}

// ToTileSet converts ts into a TileSet
func (ts TileCoordSet) ToTileSet(fn func(t TileCoord) Tile) TileSet {
	var result TileSet
	for val := range ts.set {
		result.Add(fn(val))
	}
	return result
}

func (ts TileCoordSet) String() string {
	slice := ts.Slice()

	var maxX, maxY int
	for _, tile := range slice {
		if tile.X > maxX {
			maxX = tile.X
		}
		if tile.Y > maxY {
			maxY = tile.Y
		}
	}
	maxX++
	maxY++

	tilesAt := make([][]bool, maxX)
	for x := range tilesAt {
		tilesAt[x] = make([]bool, maxY)
	}
	for _, v := range slice {
		tilesAt[v.X][v.Y] = true
	}

	var sb strings.Builder
	sb.WriteByte('{')
	for y := maxY - 1; y >= 0; y-- {
		for x := 0; x < maxX; x++ {
			if tilesAt[x][y] {
				sb.WriteByte('x')
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
func (ts TileCoordSet) MultiLineString() string {
	next := ts.String()
	next = next[1 : len(next)-1]
	next = strings.ReplaceAll(next, "|", "\n")
	next += "\n"
	return next
}
