package gridspech

import (
	"strings"
)

func stringToTileData(s string) TileData {
	var data TileData
	value := int(s[0] - '0')
	if strings.Contains(s, "_") {
		data.Type = TypeHole
	} else {
		data.Color = TileColor(value)
		data.Sticky = strings.Contains(s, "/")
		data.Type = TypeBlank
		if strings.Contains(s, "e") {
			data.Type = TypeGoal
		}
		if strings.Contains(s, "k") {
			data.Type = TypeCrown
		}
		if strings.Contains(s, "j1") {
			data.Type = TypePlus
		}
	}
	return data
}

func emptyLine(s string) bool {
	return !strings.ContainsAny(s, "_0123456789")
}

func stripEmptyLines(lines []string) []string {
	if emptyLine(lines[0]) {
		lines = lines[1:]
	}
	if emptyLine(lines[len(lines)-1]) {
		lines = lines[:len(lines)-1]
	}
	return lines
}

// MakeGridFromString takes a string and converts it into a Grid.
func MakeGridFromStringNew(str string, maxColors int) Grid {
	var grid Grid

	lines := strings.Split(str, "\n")
	lines = stripEmptyLines(lines)

	height := len(lines)
	width := len(strings.Fields(lines[0]))

	grid.Tiles = make([][]Tile, height)
	grid.MaxColors = maxColors

	for y := 0; y < height; y++ {
		grid.Tiles[y] = make([]Tile, width)
		row := strings.Fields(lines[y])

		for x := 0; x < width; x++ {
			cur := row[x]
			data := stringToTileData(cur)
			grid.Tiles[y][x] = Tile{
				Data:  data,
				Coord: TileCoord{X: x, Y: y},
			}
		}
	}

	var rotate = make([][]Tile, width)
	for x := 0; x < width; x++ {
		rotate[x] = make([]Tile, height)
		for y := 0; y < height; y++ {
			rotate[x][y] = grid.Tiles[height-1-y][x]
		}
	}
	grid.Tiles = rotate

	return grid
}

func (g Grid) String() string {
	var longest int

	tileStrs := make([][]string, g.Width())
	for i := range tileStrs {
		tileStrs[i] = make([]string, g.Height())
	}

	for y := 0; y < g.Height(); y++ {
		for x := 0; x < g.Width(); x++ {
			str := g.TileAt(x, y).Data.String()
			tileStrs[x][y] = str
			if len(str) > longest {
				longest = len(str)
			}
		}
	}

	var sb strings.Builder
	for y := 0; y < g.Height(); y++ {
		if y > 0 {
			sb.WriteByte('\n')
		}
		for x := 0; x < g.Width(); x++ {
			tileStr := tileStrs[x][y]
			padded := tileStr + strings.Repeat(" ", longest-len(tileStr))

			if x > 0 {
				sb.WriteString("  ")
			}
			sb.WriteString(padded)
		}
	}
	return sb.String()
}
