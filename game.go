package gridspech

//go:generate stringer -type=TileType -linecomment -output=stringers.go

const (
	// TypeHole represents a tile which does not exist. They cannot have Color.
	TypeHole TileType = iota // _

	// TypeBlank is a tile which does not have any icons.
	TypeBlank // Blank

	// TypeEnd is a tile which must have a direct path to another end.
	TypeEnd // End

	// TypeKing tiles must touch all tiles of their state.
	// If there are multiple crowns on the same state, they must not not touch each other, and
	// together they must touch all tiles of their state.
	TypeKing // King

	// TypeMine1 must be touching exactly 1 tiles with Color >= 1.
	TypeMine1 // Mine1

	// TypeMine2 must be touching exactly 2 tiles with Color >= 1.
	TypeMine2 // Mine2

	// TypeMine3 must be touching exactly 3 tiles with Color >= 1.
	TypeMine3 // Mine3

	// TypeJoin1 must touch exactly 1 tile with a non-blank type.
	TypeJoin1 // Join1
)

// Constants for TileColor
const (
	ColorNone TileColor = iota // _
	ColorA                     // A
	ColorB                     // B
)

// TileColor represents if a tile is disabled (0), or has a color (> 0, different colors have increasing numbers).
type TileColor byte

// TileType represents what kind of tile it is, ie "what icon to display on it".
type TileType byte

// Grid represents a single level of gridspech.
type Grid struct {
	Tiles     [][]Tile
	MaxColors int
}

// TileCoord is an X,Y coordinate in the grid.
type TileCoord struct {
	X, Y int
}

// TileData represents all of the properties which a tile can have.
type TileData struct {
	Color  TileColor
	Type   TileType
	Sticky bool

	ArrowNorth bool
	ArrowEast  bool
	ArrowSouth bool
	ArrowWest  bool
}

// Tile represents a tile in the game of gridspech. The default value of a tile will have
// `Type = Hole`.
type Tile struct {
	Coord TileCoord
	Data  TileData
}
