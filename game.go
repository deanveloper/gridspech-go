package gridspech

const (
	// TypeHole represents a tile which does not exist. They cannot have Color.
	TypeHole TileType = iota
	// TypeBlank is a tile which does not have any icons.
	TypeBlank
	// TypeGoal is a tile which must have a direct path to another goal.
	TypeGoal
	// TypeCrown tiles must touch all tiles of their state.
	// If there are multiple crowns on the same state, they must not not touch each other, and
	// together they must touch all tiles of their state.
	TypeCrown
	// TypeDot1 must be touching exactly 1 tile with Color >= 1.
	TypeDot1
	// TypeDot2 must be touching exactly 2 tile with Color >= 1.
	TypeDot2
	// TypePlus must touch exactly 1 tile with a non-blank type.
	TypePlus
)

// TileColor represents if a tile is disabled (0), or has a color (> 0, different colors have increasing numbers).
type TileColor byte

// TileType represents what kind of tile it is, ie "what icon to display on it".
type TileType byte

// Grid represents a single level of gridspech.
type Grid struct {
	Tiles [][]Tile
}

// Tile represents a tile in the game of gridspech. The default value of a tile will have
// `Type = Hole`.
type Tile struct {
	Color  TileColor
	Type   TileType
	Sticky bool
	X, Y   int
}
