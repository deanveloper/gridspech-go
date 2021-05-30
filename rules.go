package gridspech

import "fmt"

// ValidTile returns if t is valid in g. If all tiles in g are valid,
// the grid is completed.
func (g Grid) ValidTile(t Tile) bool {
	switch t.Type {
	case Hole, Blank:
		return true
	case Goal:
		return g.validGoal(t)
	case Crown:
		return g.validCrown(t)
	case Dot1:
		return g.neighborsWithState(t, Enabled).Len() == 1
	case Dot2:
		return g.neighborsWithState(t, Enabled).Len() == 2
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
		if t.Type == Goal {
			goals++

			// requirement 2: The goals should have exactly 1 neighbor with the same state.
			neighbors := g.neighborsWithState(t, t.State)
			if len(neighbors.Slice()) != 1 {
				return false
			}
		}

		// requirement 3: All other tiles in the blob should have exactly 2 neighbors with the same state.
		if t.Type != Goal && len(g.neighborsWithState(t, t.State).Slice()) != 2 {
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
		if tile.Type == Crown && tile != start {
			return false
		}
	}

	crownsWithSameState := g.TilesWith(func(t Tile) bool {
		return t.Type == Crown && t.State == start.State
	})

	// set of blobs of all crowns with same state
	var crownsBlobSet TileSet
	for _, crown := range crownsWithSameState.Slice() {
		crownsBlobSet.Merge(g.Blob(crown))
	}

	// set of all tiles with same state
	stateSet := g.TilesWith(func(t Tile) bool {
		return t.Type != Hole && t.State == start.State
	})

	// requirement 2: All tiles with the same state must have a crown in its blob.
	return crownsBlobSet.Eq(stateSet)
}
