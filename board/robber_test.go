package board

import (
	"slices"
	"testing"
)

// helper to find cross coords adjacent to a given tile
func crossesAdjacentTo(tile TileCoord) []CrossCoord {
	var result []CrossCoord
	for x := 0; x <= 5; x++ {
		for y := 0; y <= 10; y++ {
			cross, ok := NewCrossCoord(x, y)
			if !ok {
				continue
			}
			if slices.Contains(cross.adjacentTileCoords(), tile) {
				result = append(result, cross)
			}
		}
	}
	return result
}

func TestPlaceRobberReturnsAdjacentPlayerIds_Single(t *testing.T) {
	b := NewDesertBoard()

	// Choose a valid tile roughly in the middle
	tile, ok := NewTileCoord(2, 3)
	if !ok {
		t.Fatal("expected valid tile coordinate (2,3)")
	}

	adj := crossesAdjacentTo(tile)
	if len(adj) < 2 {
		t.Fatalf("expected at least 2 adjacent crosses to tile %v, got %d", tile, len(adj))
	}

	// Place one adjacent settlement for player 0 (bypass distance rule for setup)
	b.settlements[adj[0]] = 0

	// Place a far settlement for player 1 that is NOT adjacent to the tile
	farCross, _ := NewCrossCoord(4, 6)
	if !b.CanPlaceSettlement(farCross) {
		t.Fatal("expected to be able to place far settlement")
	}
	if !b.SetSettlement(farCross, 1) {
		t.Fatal("failed to set settlement for player 1")
	}

	ids := b.PlaceRobber(tile)
	if len(ids) != 1 || !slices.Contains(ids, 0) {
		t.Fatalf("expected only player 0 to be stealable, got %v", ids)
	}
}

func TestPlaceRobberReturnsAdjacentPlayerIds_Multiple(t *testing.T) {
	b := NewDesertBoard()

	tile, ok := NewTileCoord(2, 3)
	if !ok {
		t.Fatal("expected valid tile coordinate (2,3)")
	}

	adj := crossesAdjacentTo(tile)
	if len(adj) < 3 {
		t.Fatalf("expected at least 3 adjacent crosses to tile %v, got %d", tile, len(adj))
	}

	// Place two adjacent settlements for different players (bypass distance rule for setup)
	b.settlements[adj[0]] = 0
	b.settlements[adj[1]] = 1

	ids := b.PlaceRobber(tile)
	// Order is not guaranteed; verify set equality
	if len(ids) != 2 || !slices.Contains(ids, 0) || !slices.Contains(ids, 1) {
		t.Fatalf("expected players {0,1}, got %v", ids)
	}
}

func TestGetRobberReturnsCurrentPosition(t *testing.T) {
	b := NewDesertBoard()

	// Place robber on a specific tile
	tile, ok := NewTileCoord(2, 3)
	if !ok {
		t.Fatal("expected valid tile coordinate (2,3)")
	}

	b.PlaceRobber(tile)
	
	// Verify GetRobber returns the same position
	currentPos := b.GetRobber()
	if currentPos != tile {
		t.Fatalf("expected robber position %v, got %v", tile, currentPos)
	}
}
