package solve

import gs "github.com/deanveloper/gridspech-go"

// MergeSolutionsIters makes pairs of solutions from sols1 and sols2 into
// a single solution, then returns a channel of the merged pairs of solutions.
//
// A solution pair will only be sent if any tiles which appear in both solutions are equal.
func MergeSolutionsIters(sols1, sols2 <-chan gs.TileSet) <-chan gs.TileSet {
	iter := make(chan gs.TileSet, 20)

	go func() {
		// read sols2 into a slice
		var sols2slice []gs.TileSet
		for sol2 := range sols2 {
			sols2slice = append(sols2slice, sol2)
		}

		// merge
		for sol1 := range sols1 {

		nextSolution:
			for _, sol2 := range sols2slice {
				var merged gs.TileSet

				// do not merge if they have any tiles with unmatched colors
				for _, t1 := range sol1.Slice() {
					for _, t2 := range sol2.Slice() {
						if t1.Coord == t2.Coord && t1.Data != t2.Data {
							continue nextSolution
						}
					}
				}

				merged.Merge(sol1)
				merged.Merge(sol2)
				iter <- merged
			}
		}
		close(iter)
	}()

	return iter
}

func filterUnique(in <-chan gs.TileSet) <-chan gs.TileSet {
	filtered := make(chan gs.TileSet, 20)

	go func() {
		var alreadySeen []gs.TileSet
		for newSolution := range in {
			unique := true
			for _, seen := range alreadySeen {
				if newSolution.Eq(seen) {
					unique = false
					break
				}
			}
			if unique {
				alreadySeen = append(alreadySeen, newSolution)
				filtered <- newSolution
			}
		}
		close(filtered)
	}()

	return filtered
}

func filterValid(
	g GridSolver,
	tilesToValidate []gs.Tile,
	current gs.Tile,
	sols <-chan gs.TileSet,
) <-chan gs.TileSet {
	filtered := make(chan gs.TileSet, 20)

	go func() {
		defer close(filtered)
		for solution := range sols {
			newBase := g.Grid.Clone()
			newBase.ApplyTileSet(solution)

			allValid := true
			for _, tile := range tilesToValidate {
				if !newBase.ValidTile(tile.Coord) {
					allValid = false
					break
				}
			}
			if allValid {
				filtered <- solution
			}
		}
	}()

	return filtered
}

func filterHasTile(in <-chan gs.TileSet, coord gs.TileCoord) <-chan gs.TileSet {
	filtered := make(chan gs.TileSet, 20)

	go func() {
		defer close(filtered)
		for solution := range in {
			if solution.ToTileCoordSet().Has(coord) {
				filtered <- solution
			}
		}
	}()

	return filtered
}

func decorateSetBorder(g GridSolver, shapeColor gs.TileColor, tileSet gs.TileSet) <-chan gs.TileSet {
	iter := make(chan gs.TileSet)
	go func() {
		defer close(iter)

		var unknownNeighbors []gs.Tile
		for tile := range tileSet.Iter() {
			neighboringUnknowns := g.Grid.NeighborSetWith(tile.Coord, func(o gs.Tile) bool {
				return g.UnknownTiles.Has(o.Coord) &&
					!tileSet.ToTileCoordSet().Has(o.Coord)
			})
			unknownNeighbors = append(unknownNeighbors, neighboringUnknowns.Slice()...)
		}

		for permutation := range Permutation(g.Grid.MaxColors-1, len(unknownNeighbors)) {
			var setWithDecoration gs.TileSet
			setWithDecoration.Merge(tileSet)
			for i, unknown := range unknownNeighbors {
				color := permutation[i]
				if color >= int(shapeColor) {
					color++
				}
				unknown.Data.Color = gs.TileColor(color)
				setWithDecoration.Add(unknown)
			}
			iter <- setWithDecoration
		}
	}()
	return iter
}

func decorateSetIterBorders(g GridSolver, shapeColor gs.TileColor, tileSets <-chan gs.TileSet) <-chan gs.TileSet {
	iter := make(chan gs.TileSet)
	go func() {
		defer close(iter)
		for tileSet := range tileSets {

			var unknownNeighbors []gs.Tile
			for tile := range tileSet.Iter() {
				neighboringUnknowns := g.Grid.NeighborSetWith(tile.Coord, func(o gs.Tile) bool {
					return g.UnknownTiles.Has(o.Coord) &&
						!tileSet.ToTileCoordSet().Has(o.Coord)
				})
				unknownNeighbors = append(unknownNeighbors, neighboringUnknowns.Slice()...)
			}

			for permutation := range Permutation(g.Grid.MaxColors-1, len(unknownNeighbors)) {
				var setWithDecoration gs.TileSet
				setWithDecoration.Merge(tileSet)
				for i, unknown := range unknownNeighbors {
					color := permutation[i]
					if color >= int(shapeColor) {
						color++
					}
					unknown.Data.Color = gs.TileColor(color)
					setWithDecoration.Add(unknown)
				}
				iter <- setWithDecoration
			}
		}
	}()
	return iter
}

func mergeSolutionsSlices(sols1, sols2 []gs.TileSet) []gs.TileSet {
	var result []gs.TileSet
	for _, sol1 := range sols1 {
	nextSolution:
		for _, sol2 := range sols2 {

			// do not merge if they have any tiles with unmatched colors
			for _, t1 := range sol1.Slice() {
				for _, t2 := range sol2.Slice() {
					if t1.Coord == t2.Coord && t1.Data.Color != t2.Data.Color {
						continue nextSolution
					}
				}
			}

			var merged gs.TileSet
			merged.Merge(sol1)
			merged.Merge(sol2)
			result = append(result, merged)
		}
	}
	return result
}

func removeIfInvalid(g GridSolver, tilesToValidate []gs.TileCoord, in []gs.TileSet) []gs.TileSet {
	var validSolutions []gs.TileSet

	base := g.Grid
	for _, solution := range in {
		newBase := base.Clone()
		newBase.ApplyTileSet(solution)

		allValid := true
		for _, coord := range tilesToValidate {
			if !newBase.ValidTile(coord) {
				allValid = false
				break
			}
		}
		if allValid {
			validSolutions = append(validSolutions, solution)
		}
	}

	return validSolutions
}

func removeIfNonUnique(in []gs.TileSet) []gs.TileSet {
	var filtered []gs.TileSet

	for _, solution := range in {
		unique := true
		for _, seen := range filtered {
			if solution.Eq(seen) {
				unique = false
				break
			}
		}
		if unique {
			filtered = append(filtered, solution)
		}
	}

	return filtered
}
