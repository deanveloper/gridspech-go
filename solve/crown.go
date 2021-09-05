package solve

import (
	"fmt"

	gs "github.com/deanveloper/gridspech-go"
)

// Crowns will return a channel of solutions for all the crown tiles in g.
func Crowns(g GridSolver) <-chan gs.TileSet {
	return nil
}

// crowns are very tough to solve, so we basically just have to "abort" if any tiles which are not crowns become invalid.
func (g GridSolver) solveCrown() {
}

// SolveShapes returns an iterator of all shapes which contain `start`, and be made out `color`, as well
// as a communication channel to say whether this branch should be traversed or not.
// The two channels will be closed after
func (g GridSolver) SolveShapes(start gs.TileCoord, color gs.TileColor) (<-chan gs.TileSet, chan<- bool) {
	solutionsChan := make(chan gs.TileSet)
	pruneChan := make(chan bool)

	go func() {
		bfsShapes(g, start, color, solutionsChan, pruneChan)
		close(solutionsChan)
		close(pruneChan)
	}()

	return solutionsChan, pruneChan
}

func bfsShapes(g GridSolver, start gs.TileCoord, color gs.TileColor, solutions chan<- gs.TileSet, pruneChan <-chan bool) {
	var blobQueue []gs.TileCoordSet
	blobQueue = append(blobQueue, gs.NewTileCoordSet(start))

	var dupeChecker []gs.TileCoordSet
	blobSize := 1

	for len(blobQueue) > 0 {

		curShape := blobQueue[0]
		blobQueue = blobQueue[1:]

		if curShape.Len() > blobSize {
			dupeChecker = nil
			blobSize = curShape.Len()
		}

		tileSet := curShape.ToTileSet(func(t gs.TileCoord) gs.Tile {
			tileCopy := *g.Grid.TileAtCoord(t)
			tileCopy.Data.Color = color
			return tileCopy
		})
		solutions <- tileSet
		if <-pruneChan {
			continue
		}

		var allNeighbors gs.TileCoordSet
		for _, tile := range curShape.Slice() {
			newNeighbors := g.Grid.NeighborsWith(tile, func(o gs.Tile) bool {
				return !curShape.Has(o.Coord) && (g.UnknownTiles.Has(o.Coord) || o.Data.Color == color)
			})
			allNeighbors.Merge(newNeighbors.ToTileCoordSet())
		}

	neighborLoop:
		for _, newNeighbor := range allNeighbors.Slice() {
			var newShape gs.TileCoordSet
			newShape.Merge(curShape)
			newShape.Add(newNeighbor)

			// check if newShape has already been done
			for _, dupe := range dupeChecker {
				if dupe.Eq(newShape) {
					continue neighborLoop
				}
			}

			blobQueue = append(blobQueue, newShape)
			dupeChecker = append(dupeChecker, newShape)
			fmt.Println(dupeChecker)
		}
	}
}
