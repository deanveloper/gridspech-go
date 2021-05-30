package gridspech

// NorthOf returns the tile north of t in g.
func (g Grid) NorthOf(t Tile) Tile {
	if t.Y == g.Height-1 || t.Type == Hole {
		return Tile{}
	}
	return g.Tiles[t.X][t.Y+1]
}

// EastOf returns the tile east of t in g.
func (g Grid) EastOf(t Tile) Tile {
	if t.X == g.Width-1 || t.Type == Hole {
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
	neighbors := g.neighborsWithState(t, t.State)

	for _, neighbor := range neighbors.Slice() {
		if !ts.Has(neighbor) {
			ts.Add(neighbor)
			g.blobRecur(neighbor, ts)
		}
	}
}

// returns the set of neighbors which have a certain state
func (g Grid) neighborsWithState(t Tile, state TileState) TileSet {
	neighbors := g.Neighbors(t)
	for _, neighbor := range neighbors.Slice() {
		if neighbor.State != state {
			neighbors.Remove(neighbor)
		}
	}
	return neighbors
}
