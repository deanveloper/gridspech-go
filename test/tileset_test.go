package gridspech_test

import (
	"testing"

	gs "github.com/deanveloper/gridspech-go"
)

func TestTileSetAddRemove(t *testing.T) {
	var ts gs.TileSet
	if ts.Len() != 0 {
		t.Errorf("Length test failed. Expected %v, got %v", 0, ts.Len())
	}
	tile := gs.Tile{Data: gs.TileData{Color: 2}}
	ts.Add(tile)
	if ts.Len() != 1 {
		t.Errorf("Length test after adding the first time failed. Expected %v, got %v", 1, ts.Len())
	}
	if !ts.Has(tile) {
		t.Errorf("expected ts to have %#v after adding, but it did not.", tile)
	}
	ts.Add(tile)
	if ts.Len() != 1 {
		t.Errorf("Length test after adding a second time failed. Expected %v, got %v", 1, ts.Len())
	}
	ts.Remove(tile)
	if ts.Has(tile) {
		t.Errorf("expected ts not to have %#v after removing, but it was still in ts.", tile)
	}
	if ts.Len() != 0 {
		t.Errorf("Length test after removing failed. Expected %v, got %v", 0, ts.Len())
	}
}

func TestTileSetHas(t *testing.T) {
	ts := gs.NewTileSet(gs.Tile{Data: gs.TileData{Color: 2}}, gs.Tile{Data: gs.TileData{Color: 10}}, gs.Tile{Data: gs.TileData{Color: 4, Type: gs.TypeCrown}})

	cases := []struct {
		Value    gs.Tile
		Expected bool
	}{
		{gs.Tile{Data: gs.TileData{Color: 2}}, true},
		{gs.Tile{Data: gs.TileData{Color: 10}}, true},
		{gs.Tile{Data: gs.TileData{Color: 4, Type: gs.TypeCrown}}, true},
		{gs.Tile{}, false},
		{gs.Tile{Data: gs.TileData{Type: gs.TypeCrown}}, false},
	}

	for _, testCase := range cases {
		if ts.Has(testCase.Value) != testCase.Expected {
			t.Errorf("failed\nexpected ts.Has(%#v) to be %v, but was %v", testCase.Value, testCase.Expected, ts.Has(testCase.Value))
		}
	}
}
