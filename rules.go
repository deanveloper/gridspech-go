package gridspech

import "fmt"

// ValidTile returns if t is valid in g. If all tiles in g are valid,
// the grid is completed.
func (g Grid) ValidTile(coord TileCoord) bool {
	t := g.TileAtCoord(coord)

	switch t.Data.Type {
	case TypeHole, TypeBlank:
		return true
	case TypeGoal:
		return g.validGoal(*t)
	case TypeCrown:
		return g.validCrown(*t)
	case TypeDot1:
		return g.NeighborsWith(t.Coord, func(other Tile) bool {
			return other.Data.Color != ColorNone && other.Data.Color != 100
		}).Len() == 1
	case TypeDot2:
		return g.NeighborsWith(t.Coord, func(other Tile) bool {
			return other.Data.Color != ColorNone && other.Data.Color != 100
		}).Len() == 2
	case TypeDot3:
		return g.NeighborsWith(t.Coord, func(other Tile) bool {
			return other.Data.Color != ColorNone && other.Data.Color != 100
		}).Len() == 3
	case TypePlus:
		return g.validPlus(*t)
	default:
		panic(fmt.Sprintf("invalid tile type %v", t.Data.Type))
	}
}

// the blob of a goal tile should contain a direct path to another goal.
// the way we measure this:
//   1. The blob should contain exactly two goals.
//   2. The goals should have exactly 1 neighbor with the same state.
//   3. All other tiles in the blob should have exactly 2 neighbors with the same state.
func (g Grid) validGoal(start Tile) bool {
	blob := g.Blob(start.Coord)
	var goals int
	for _, t := range blob.Slice() {
		if t.Data.Type == TypeGoal {
			goals++

			// requirement 2: The goals should have exactly 1 neighbor with the same state.
			neighbors := g.NeighborsWith(t.Coord, func(o Tile) bool {
				return t.Data.Color == o.Data.Color
			})
			if len(neighbors.Slice()) != 1 {
				return false
			}
		}

		// requirement 3: All other tiles in the blob should have exactly 2 neighbors with the same state.
		neighborsSameColor := g.NeighborsWith(t.Coord, func(o Tile) bool {
			return t.Data.Color == o.Data.Color
		}).Slice()
		if t.Data.Type != TypeGoal && len(neighborsSameColor) != 2 {
			return false
		}
	}

	// requirement 1: The blob should contain exactly two goals.
	return goals == 2
}

// crown tiles have the following requirements:
//   1. No other crowns may be in this crown's blob.
//   2. All tiles with the same state must have a crown in its blob.
func (g Grid) validCrown(start Tile) bool {
	blob := g.Blob(start.Coord)

	// requirement 1: No other crowns may be in this crown's blob.
	for _, tile := range blob.Slice() {
		if tile.Data.Type == TypeCrown && tile != start {
			return false
		}
	}

	crownsWithSameState := g.TilesWith(func(t Tile) bool {
		return t.Data.Type == TypeCrown && t.Data.Color == start.Data.Color
	})

	// set of blobs of all crowns with same state
	var crownsBlobSet TileSet
	for crown := range crownsWithSameState.Iter() {
		crownsBlobSet.Merge(g.Blob(crown.Coord))
	}

	// set of all tiles with same state
	stateSet := g.TilesWith(func(t Tile) bool {
		return t.Data.Type != TypeHole && t.Data.Color == start.Data.Color
	})

	// requirement 2: All tiles with the same state must have a crown in its blob.
	return crownsBlobSet.Eq(stateSet)
}

func (g Grid) validPlus(t Tile) bool {
	var foundExactlyOne bool
	for _, blobTile := range g.Blob(t.Coord).Slice() {
		if blobTile.Data.Type != TypeHole && blobTile.Data.Type != TypeBlank {
			if foundExactlyOne {
				return false
			}
			foundExactlyOne = true
		}
	}
	return foundExactlyOne
}
