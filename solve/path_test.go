package solve_test

import (
	"testing"

	gs "github.com/deanveloper/gridspech-go"
	"github.com/deanveloper/gridspech-go/example"
	"github.com/deanveloper/gridspech-go/solve"
)

func testPathsIterAbstract(t *testing.T, level string, x1, y1, x2, y2 int, solutions []string) {
	t.Helper()

	grid := solve.NewGridSolver(gs.MakeGridFromString(level, 2))

	var expectedSolutions []gs.TileSet
	for i := range solutions {
		expectedSolutions = append(expectedSolutions, tileSetFromString(grid.Grid, solutions[i]))
	}

	ch := grid.PathsIter(gs.TileCoord{X: x1, Y: y1}, gs.TileCoord{X: x2, Y: y2}, 1)
	var actualSolutions []gs.TileSet
	for ts := range ch {
		actualSolutions = append(actualSolutions, ts)
	}

	testUnorderedTilesetSliceEq(t, expectedSolutions, actualSolutions)
}

func TestPathsIter_levelA1(t *testing.T) {
	const level = example.LevelA1

	solutions := []string{"1111"}
	testPathsIterAbstract(t, level, 0, 0, 3, 0, solutions)
}

func TestPathsIter_levelA2(t *testing.T) {
	const level = example.LevelA2
	solutions := []string{"1 11|111 "}
	testPathsIterAbstract(t, level, 0, 1, 3, 1, solutions)
}

func TestPathsIter_levelA3(t *testing.T) {
	const level = example.LevelA3
	solutions := []string{
		" 111 |11 11|     ",
		"     |11 11| 111 ",
	}
	testPathsIterAbstract(t, level, 0, 1, 4, 1, solutions)
}

func TestPathsIter_levelA4(t *testing.T) {
	const level = example.LevelA4
	solutions := []string{"1 1|111"}
	testPathsIterAbstract(t, level, 0, 1, 2, 1, solutions)
}

func TestPathsIter_levelA5(t *testing.T) {
	const level = example.LevelA5
	solutions := []string{"111 |1 11|1  1"}
	testPathsIterAbstract(t, level, 0, 0, 3, 0, solutions)
}

func TestPathsIter_levelA6(t *testing.T) {
	const level = example.LevelA6
	solutions := []string{
		"111 111 |1 1 1 11|1 111  1",
		"111 111 |1 1 1 1 |1 111 11",
		"    111 |111 1 11|1 111  1",
		"    111 |111 1 1 |1 111 11",
	}
	testPathsIterAbstract(t, level, 0, 0, 7, 0, solutions)
}

func TestPathsIter_levelA9(t *testing.T) {
	const level = example.LevelA9
	solutions := []string{
		"   1111|1 11  1|111    ",
		"   111 |1 11 11|111    ",
	}
	testPathsIterAbstract(t, level, 0, 1, 6, 1, solutions)
}

func TestPathsIter_basicColorNonePath(t *testing.T) {
	const level = `
	1/   0    0    1/ 
	0/e  0    0    0e 
	1/   0    0    1/ 
	`
	grid := gs.MakeGridFromString(level, 2)
	gridSolver := solve.NewGridSolver(grid)
	solutionsChan := gridSolver.PathsIter(gs.TileCoord{X: 0, Y: 1}, gs.TileCoord{X: 3, Y: 1}, gs.ColorNone)
	var solutions []gs.TileSet
	for solution := range solutionsChan {
		solutions = append(solutions, solution)
	}
	if len(solutions) != 1 {
		t.Fatalf("solutions length expected to be 1 but was %d", len(solutions))
	}
	expected := gs.NewTileSet(
		gs.Tile{Coord: gs.TileCoord{X: 0, Y: 1}, Data: gs.TileData{Sticky: true, Type: gs.TypeGoal}},
		gs.Tile{Coord: gs.TileCoord{X: 1, Y: 1}, Data: gs.TileData{Type: gs.TypeBlank}},
		gs.Tile{Coord: gs.TileCoord{X: 2, Y: 1}, Data: gs.TileData{Type: gs.TypeBlank}},
		gs.Tile{Coord: gs.TileCoord{X: 3, Y: 1}, Data: gs.TileData{Type: gs.TypeGoal}},
	)

	if !expected.Eq(solutions[0]) {
		expectedGrid := grid.Clone()
		actualGrid := grid.Clone()
		expectedGrid.ApplyTileSet(expected)
		actualGrid.ApplyTileSet(solutions[0])
		t.Errorf("solutions not equal")
		t.Errorf("expected:\n%v", expectedGrid)
		t.Errorf("actual:\n%v", actualGrid)
	}
}
