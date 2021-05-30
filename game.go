package gridspech

const (
	// Disabled is the disabled tile state.
	Disabled TileState = iota
	// Enabled is the Enabled tile state.
	Enabled
)

const (
	// Hole represents a tile which does not exist. They cannot have state.
	Hole TileType = iota
	// Blank is a tile which does not have any icons.
	Blank
	// Goal is a tile which must have a direct path to another goal.
	Goal
	// Crown tiles must touch all tiles of their state.
	// If there are multiple crowns on the same state, they must not not touch each other, and
	// together they must touch all tiles of their state.
	Crown
	// Dot1 must be touching exactly 1 enabled tile.
	Dot1
	// Dot2 must be touching exactly 2 enabled tiles.
	Dot2
)

// TileState represents if a tile is Enabled or Disabled.
type TileState byte

// TileType represents what kind of tile it is, ie "what icon to display on it".
type TileType byte

// Grid represents a game.
type Grid struct {
	Tiles         [][]Tile
	Width, Height int
}

// Tile represents a tile in the game of gridspech. The default value of a tile will have
// `Type = Hole`.
type Tile struct {
	State  TileState
	Type   TileType
	Sticky bool
	X, Y   int
}
