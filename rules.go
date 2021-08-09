package gridspech

import "fmt"

// ValidTile returns if t is valid in g. If all tiles in g are valid,
// the grid is completed.
func (g Grid) ValidTile(t Tile) bool {
	switch t.Type {
	case TypeHole, TypeBlank:
		return true
	case TypeGoal:
		return g.validGoal(t)
	case TypeCrown:
		return g.validCrown(t)
	case TypeDot1:
		return g.NeighborsWith(t, func(other Tile) bool {
			return other.Color != ColorNone && other.Color != 100
		}).Len() == 1
	case TypeDot2:
		return g.NeighborsWith(t, func(other Tile) bool {
			return other.Color != ColorNone && other.Color != 100
		}).Len() == 2
	case TypeDot3:
		return g.NeighborsWith(t, func(other Tile) bool {
			return other.Color != ColorNone && other.Color != 100
		}).Len() == 3
	case TypePlus:
		return g.validPlus(t)
	default:
		panic(fmt.Sprintf("invalid tile type %v", t.Type))
	}
}

// the blob of a goal tile should contain a direct path to another goal.
// the way we measure this:
//   1. The blob should contain exactly two goals.
//   2. The goals should have exactly 1 neighbor with the same state.
//   3. All other tiles in the blob should have exactly 2 neighbors with the same state.
func (g Grid) validGoal(start Tile) bool {
	blob := g.Blob(start)
	var goals int
	for _, t := range blob.Slice() {
		if t.Type == TypeGoal {
			goals++

			// requirement 2: The goals should have exactly 1 neighbor with the same state.
			neighbors := g.NeighborsWith(t, func(o Tile) bool {
				return t.Color == o.Color
			})
			if len(neighbors.Slice()) != 1 {
				return false
			}
		}

		// requirement 3: All other tiles in the blob should have exactly 2 neighbors with the same state.
		neighborsSameColor := g.NeighborsWith(t, func(o Tile) bool {
			return t.Color == o.Color
		}).Slice()
		if t.Type != TypeGoal && len(neighborsSameColor) != 2 {
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
	blob := g.Blob(start)

	// requirement 1: No other crowns may be in this crown's blob.
	for _, tile := range blob.Slice() {
		if tile.Type == TypeCrown && tile != start {
			return false
		}
	}

	crownsWithSameState := g.TilesWith(func(t Tile) bool {
		return t.Type == TypeCrown && t.Color == start.Color
	})

	// set of blobs of all crowns with same state
	var crownsBlobSet TileSet
	for _, crown := range crownsWithSameState.Slice() {
		crownsBlobSet.Merge(g.Blob(crown))
	}

	// set of all tiles with same state
	stateSet := g.TilesWith(func(t Tile) bool {
		return t.Type != TypeHole && t.Color == start.Color
	})

	// requirement 2: All tiles with the same state must have a crown in its blob.
	return crownsBlobSet.Eq(stateSet)
}

func (g Grid) validPlus(t Tile) bool {
	var foundExactlyOne bool
	for _, blobTile := range g.Blob(t).Slice() {
		if blobTile.Type != TypeHole && blobTile.Type != TypeBlank {
			if foundExactlyOne {
				return false
			}
			foundExactlyOne = true
		}
	}
	return foundExactlyOne
}
