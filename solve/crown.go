package solve

import (
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
	type shape struct {
		fullSet, newTiles gs.TileCoordSet
	}

	var blobQueue []shape
	blobQueue = append(blobQueue, shape{gs.NewTileCoordSet(start), gs.NewTileCoordSet(start)})

	for len(blobQueue) > 0 {

		curShape := blobQueue[0]
		blobQueue = blobQueue[1:]

		tileSet := curShape.fullSet.ToTileSet(func(t gs.TileCoord) gs.Tile {
			tileCopy := *g.Grid.TileAtCoord(t)
			tileCopy.Data.Color = color
			return tileCopy
		})
		solutions <- tileSet
		if <-pruneChan {
			continue
		}

		var allNewNeighbors gs.TileCoordSet
		for _, newTile := range curShape.newTiles.Slice() {
			newTileNeighbors := g.Grid.NeighborsWith(newTile, func(o gs.Tile) bool {
				return !curShape.fullSet.Has(o.Coord) && (g.UnknownTiles.Has(o.Coord) || o.Data.Color == color)
			})
			allNewNeighbors.Merge(newTileNeighbors.ToTileCoordSet())
		}
		allNewNeighborsSlice := allNewNeighbors.Slice()
		for perm := range Permutation(2, len(allNewNeighborsSlice)) {
			var newNeighbors gs.TileCoordSet
			for i := range perm {
				if perm[i] == 1 {
					newNeighbors.Add(allNewNeighborsSlice[i])
				}
			}
			if newNeighbors.Len() == 0 {
				continue
			}
			var newShape gs.TileCoordSet
			newShape.Merge(curShape.fullSet)
			newShape.Merge(newNeighbors)
			blobQueue = append(blobQueue, shape{fullSet: newShape, newTiles: newNeighbors})
		}
	}
}
