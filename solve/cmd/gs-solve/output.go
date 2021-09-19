package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/deanveloper/gridspech-go"
	"github.com/deanveloper/gridspech-go/solve"
)

func printSolutionAsLines(solver solve.GridSolver, solutions <-chan gridspech.TileSet) {
	first := true
	for solution := range solutions {
		if !first {
			fmt.Println()
		}
		first = false

		newGrid := solver.Grid.Clone()
		newGrid.ApplyTileSet(solution)
		fmt.Println(newGrid)
	}
}

func printSolutionAsJSON(solver solve.GridSolver, solutions <-chan gridspech.TileSet) {
	var solutionsArray [][][]gridspech.Tile
	for solution := range solutions {
		newGrid := solver.Grid.Clone()
		newGrid.ApplyTileSet(solution)
		solutionsArray = append(solutionsArray, newGrid.Tiles)
	}
	json.NewEncoder(os.Stdout).Encode(solutionsArray)
}

func printTileSetAsEmoji(w io.Writer, width, height int, solution map[gridspech.TileCoord]gridspech.Tile) {
	emojis := []string{
		"⬛",
		"🟥",
		"🟨",
		"🟦",
		"🟩",
		"🟪",
		"🟫",
		"⬜",
		"🟧",
	}
	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			tile, ok := solution[gridspech.TileCoord{X: x, Y: y}]
			if ok {
				fmt.Fprint(w, emojis[tile.Data.Color])
			} else {
				fmt.Fprint(w, emojis[0])
			}
		}
		fmt.Fprintln(w)
	}
}

func printSolutionAsBlocks(solver solve.GridSolver, solutions <-chan gridspech.TileSet) {
	first := true

	width, height := solver.Grid.Width(), solver.Grid.Height()
	for solution := range solutions {
		if !first {
			fmt.Println()
		}
		first = false

		solutionMap := make(map[gridspech.TileCoord]gridspech.Tile)
		for tile := range solution.Iter() {
			solutionMap[tile.Coord] = tile
		}

		printTileSetAsEmoji(os.Stdout, width, height, solutionMap)
	}
}
