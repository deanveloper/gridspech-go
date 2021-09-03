package solve

import gs "github.com/deanveloper/gridspech-go"

// Kings will return a channel of solutions for all the king tiles in g.
func Kings(g GridSolver) <-chan gs.TileSet {

	return nil
}

// crowns are very tough to solve, so we basically just have to "abort" if any tiles which are not crowns become invalid.
func (g GridSolver) solveKings() <-chan gs.TileSet {
	kingTiles := g.Grid.TilesWith(func(o gs.Tile) bool {
		return o.Data.Type == gs.TypeKing
	}).Slice()

	_ = kingTiles
	return nil
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
	var queue []gs.TileCoord
	queue = append(queue, start)
	remainingInLayer := 1

	var dupeChecker []gs.TileCoordSet

	var currentSet gs.TileCoordSet
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		remainingInLayer--

		neighbors := g.Grid.NeighborsWith(current, func(o gs.Tile) bool {
			return (o.Data.Color == color || g.UnknownTiles.Has(o.Coord)) && !currentSet.Has(o.Coord)
		})

	nextLoop:
		for _, next := range neighbors.Slice() {
			currentSet.Add(next.Coord)
			for _, prevSet := range dupeChecker {
				if prevSet.Eq(currentSet) {
					currentSet.Remove(next.Coord)
					continue nextLoop
				}
			}
			var newSet gs.TileCoordSet
			newSet.Merge(currentSet)
			dupeChecker = append(dupeChecker, newSet)
			tileSet := newSet.ToTileSet(func(t gs.TileCoord) gs.Tile { return *g.Grid.TileAtCoord(t) })
			solutions <- tileSet
			if !<-pruneChan {
				queue = append(queue, next.Coord)
			}
		}

		if remainingInLayer == 0 {
			dupeChecker = nil
			remainingInLayer = len(queue)
		}
	}
}
