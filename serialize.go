package gridspech

import (
	"regexp"
	"strings"
)

func deserializeTile(s string, x, y int) Tile {
	var tile Tile
	value := int(s[0] - '0')
	if strings.Contains(s, "_") {
		tile.Data.Type = TypeHole
	} else {
		tile.Data.Color = TileColor(value)
		tile.Data.Sticky = strings.Contains(s, "/")
		tile.Data.Type = TypeBlank
		if strings.Contains(s, "e") {
			tile.Data.Type = TypeGoal
		}
		if strings.Contains(s, "k") {
			tile.Data.Type = TypeCrown
		}
		if strings.Contains(s, "j1") {
			tile.Data.Type = TypePlus
		}
	}
	tile.Coord.X = x
	tile.Coord.Y = y
	return tile
}

func emptyLine(s string) bool {
	return !regexp.MustCompile(".*[_0-9].*").Match([]byte(s))
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

func DeserializeGrid(str string, maxColors int) Grid {
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
		println(lines[y])

		for x := 0; x < width; x++ {
			cur := row[x]
			tile := deserializeTile(cur, x, y)
			grid.Tiles[y][x] = tile
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

	println(grid.String())

	return grid
}
