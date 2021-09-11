package solve

import (
	"container/heap"

	gs "github.com/deanveloper/gridspech-go"
)

// ===== heap structure for "container/heap" =====

type blobHeap []gs.TileCoordSet

var _ heap.Interface = &blobHeap{}

func (e *blobHeap) Push(v interface{}) {
	*e = append(*e, v.(gs.TileCoordSet))
}

func (e *blobHeap) Pop() interface{} {
	elem := (*e)[len(*e)-1]
	*e = (*e)[:len(*e)-1]
	return elem
}

func (e blobHeap) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e blobHeap) Less(i, j int) bool {
	return e[i].Len() < e[j].Len()
}

func (e blobHeap) Len() int {
	return len(e)
}

// ShapesIter returns an iterator of all shapes which contain `start`, and be made out `color`, as well
// as a communication channel to say whether we should prune the set here or not.
// Both returned channels are closed by this function.
func (g GridSolver) ShapesIter(start gs.TileCoord, color gs.TileColor) (<-chan gs.TileSet, chan<- bool) {
	solutionsChan := make(chan gs.TileSet)
	pruneChan := make(chan bool)

	go func() {
		defer close(solutionsChan)
		defer close(pruneChan)

		g.bfsShapes(start, color, solutionsChan, pruneChan)
	}()

	return solutionsChan, pruneChan
}

func (g GridSolver) bfsShapes(start gs.TileCoord, color gs.TileColor, solutions chan<- gs.TileSet, pruneChan <-chan bool) {

	initialBlob := g.Grid.BlobWith(start, func(o gs.Tile) bool {
		return o.Data.Color == color && !g.UnknownTiles.Has(o.Coord)
	}).ToTileCoordSet()
	initialBlob.Add(start)

	var blobPQ blobHeap
	heap.Init(&blobPQ)
	heap.Push(&blobPQ, initialBlob)

	var dupeChecker []gs.TileCoordSet
	blobSize := 1

	for blobPQ.Len() > 0 {

		curShape := heap.Pop(&blobPQ).(gs.TileCoordSet)

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

		nextNeighbors := g.aroundShape(curShape, func(o gs.Tile) bool {
			return g.UnknownTiles.Has(o.Coord) || o.Data.Color == color
		})

	neighborLoop:
		for _, nextNeighbor := range nextNeighbors.Slice() {
			var newShape gs.TileCoordSet
			newShape.Merge(curShape)
			newShape.Add(nextNeighbor)

			// special behavior: if nextNeighbor has any neighbors which we know are the same color,
			// add the blob of each of those neighbors to ne`wShape
			transitiveNeighbors := g.Grid.NeighborsWith(nextNeighbor, func(o gs.Tile) bool {
				return o.Data.Color == color && !g.UnknownTiles.Has(o.Coord)
			})
			for _, transitiveNeighbor := range transitiveNeighbors.Slice() {
				neighborBlob := g.Grid.BlobWith(transitiveNeighbor.Coord, func(o gs.Tile) bool {
					return !g.UnknownTiles.Has(o.Coord)
				})
				newShape.Merge(neighborBlob.ToTileCoordSet())
			}

			// check if newShape has already been done
			for _, dupe := range dupeChecker {
				if dupe.Eq(newShape) {
					continue neighborLoop
				}
			}

			heap.Push(&blobPQ, newShape)
			dupeChecker = append(dupeChecker, newShape)
		}
	}
}

func (g GridSolver) aroundShape(shape gs.TileCoordSet, filter func(o gs.Tile) bool) gs.TileCoordSet {
	var allNeighbors gs.TileCoordSet
	for _, tile := range shape.Slice() {
		newNeighbors := g.Grid.NeighborsWith(tile, func(o gs.Tile) bool {
			return !shape.Has(o.Coord) && filter(o)
		})
		allNeighbors.Merge(newNeighbors.ToTileCoordSet())
	}

	return allNeighbors
}
