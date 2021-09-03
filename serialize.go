package gridspech

import (
	"strings"
)

func stringToTileData(s string) TileData {
	var data TileData
	if strings.Contains(s, "_") {
		data.Type = TypeHole
	} else {
		data.Color = TileColor(s[0] - '0')
		data.Sticky = strings.Contains(s, "/")
		data.Type = TypeBlank
		switch {
		case strings.Contains(s, "k"):
			data.Type = TypeKing
		case strings.Contains(s, "e"):
			data.Type = TypeEnd
		case strings.Contains(s, "m1"):
			data.Type = TypeMine1
		case strings.Contains(s, "m2"):
			data.Type = TypeMine2
		case strings.Contains(s, "m3"):
			data.Type = TypeMine3
		case strings.Contains(s, "j1"):
			data.Type = TypeJoin1
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
func MakeGridFromString(str string, maxColors int) Grid {
	var grid Grid

	lines := strings.Split(str, "\n")
	lines = stripEmptyLines(lines)

	height := len(lines)
	width := len(strings.Fields(lines[0]))

	grid.Tiles = make([][]Tile, width)
	for x := range grid.Tiles {
		grid.Tiles[x] = make([]Tile, height)
	}
	grid.MaxColors = maxColors

	for y := 0; y < height; y++ {
		row := strings.Fields(lines[height-y-1])

		for x := 0; x < width; x++ {
			cur := row[x]
			data := stringToTileData(cur)
			grid.Tiles[x][y] = Tile{
				Data:  data,
				Coord: TileCoord{X: x, Y: y},
			}
		}
	}

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
	for y := g.Height() - 1; y >= 0; y-- {
		if y < g.Height()-1 {
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
