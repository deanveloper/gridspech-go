package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/deanveloper/gridspech-go"
	"github.com/deanveloper/gridspech-go/solve"
	"github.com/pborman/getopt/v2"
)

const (
	outputString = "lines"
	outputJSON   = "json"
)

var (
	helpFlag    = getopt.BoolLong("help", 'h', "display help")
	maxColors   = getopt.IntLong("maxcolors", 'm', 0, "the total number of colors available for this level", "2")
	solveTiles  = getopt.ListLong("tiles", 't', "solve specific tiles. a comma-separated list of space-separated coordinates")
	solveGoals  = getopt.BoolLong("goals", 'g', "solve all goal tiles")
	solveCrowns = getopt.BoolLong("crowns", 'c', "solve all crown tiles")
	solveDots   = getopt.BoolLong("dots", 'd', "solve all dot tiles")
	solveJoins  = getopt.BoolLong("joins", 'j', "solve all join tiles")
	solveAll    = getopt.BoolLong("all", 'a', "solve all tiles")
	jsonOutput  = getopt.EnumLong("format", 'f', []string{outputString, outputJSON}, "", "output format (lines or json)")
)

func solutionsFromFlags(solver solve.GridSolver) <-chan gridspech.TileSet {
	var ch <-chan gridspech.TileSet
	if *solveAll {
		ch = solver.SolveAllTiles()
	} else {
		{
			tempCh := make(chan gridspech.TileSet, 1)
			tempCh <- gridspech.NewTileSet()
			close(tempCh)
			ch = tempCh
		}

		if getopt.IsSet('t') {
			tiles := parseCoords(*solveTiles)
			ch = solve.MergeSolutionsIters(ch, solver.SolveTiles(tiles...))
		}
		if *solveGoals {
			ch = solve.MergeSolutionsIters(ch, solver.SolveGoals())
		}
		if *solveJoins {
			ch = solve.MergeSolutionsIters(ch, solver.SolveGoals())
		}
		if *solveDots {
			ch = solve.MergeSolutionsIters(ch, solver.SolveGoals())
		}
		if *solveCrowns {
			ch = solve.MergeSolutionsIters(ch, solver.SolveGoals())
		}
	}
	return ch
}

func parseCoords(coordsStr []string) []gridspech.TileCoord {
	var coords []gridspech.TileCoord
	for _, coordStr := range coordsStr {
		var x, y int
		n, err := fmt.Sscanf(coordStr, "%d %d", &x, &y)
		if err != nil || n != 2 {
			log.Printf("skipping invalid coord %s\n", coordStr)
			continue
		}
		coords = append(coords, gridspech.TileCoord{X: x, Y: y})
	}
	return coords
}

func main() {
	getopt.HelpColumn = 22
	getopt.SetUsage(func() {
		fmt.Fprintf(
			os.Stderr, "Usage: %v %v\n",
			getopt.CommandLine.Program(),
			getopt.CommandLine.UsageLine(),
		)
		fmt.Fprintln(os.Stderr, "Standard input will be interpreted as the level to solve.")
		getopt.CommandLine.PrintOptions(os.Stderr)
	})
	getopt.Parse()
	if !getopt.IsSet('a') && !getopt.IsSet('t') && !getopt.IsSet('g') && !getopt.IsSet('c') && !getopt.IsSet('d') && !getopt.IsSet('j') {
		getopt.Usage()
		return
	}
	if *helpFlag {
		getopt.Usage()
		return
	}

	const maxLevelLen, bufLen = 10000, 100
	var buf [bufLen]byte
	var levelBytes []byte

	for i := 0; i < maxLevelLen/bufLen; i++ {
		n, err := os.Stdin.Read(buf[:])
		if err != nil {
			if err != io.EOF {
				log.Fatalln("error:", err)
			} else {
				break
			}
		}
		levelBytes = append(levelBytes, buf[:n]...)
	}
	if len(levelBytes) == maxLevelLen {
		log.Fatalln("standard input is over 10000 bytes... are you sure it is a gridspech level?")
	}

	level := string(levelBytes)
	solver := solve.NewGridSolver(gridspech.MakeGridFromString(level, *maxColors))

	solutions := solutionsFromFlags(solver)

	first := true
	for solution := range solutions {
		if !first {
			fmt.Println()
		}

		newGrid := solver.Grid.Clone()
		newGrid.ApplyTileSet(solution)
		fmt.Println(newGrid)
	}
}
