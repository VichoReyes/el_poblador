package board

import (
	"testing"
)

func TestRobberBlocksResourceGeneration(t *testing.T) {
	b := NewDesertBoard()
	
	// Set up a specific tile with wheat and dice number 6
	wheatTileCoord := TileCoord{X: 2, Y: 3}
	b.tiles[wheatTileCoord] = Tile{Terrain: TerrainWheat, DiceNumber: 6}
	
	// Place a settlement at a cross adjacent to the wheat tile
	cross := CrossCoord{X: 2, Y: 4}
	b.settlements[cross] = 0 // player 0
	
	// Generate resources without robber
	resourcesBefore := b.GenerateResources(6)
	player0ResourcesBefore := len(resourcesBefore[0])
	
	if player0ResourcesBefore == 0 {
		t.Fatal("Settlement should generate resources from adjacent wheat tile")
	}
	
	// Place robber on the wheat tile
	b.PlaceRobber(wheatTileCoord)
	
	// Generate resources with robber on the tile
	resourcesAfter := b.GenerateResources(6)
	player0ResourcesAfter := len(resourcesAfter[0])
	
	// Verify robber blocks resources
	if player0ResourcesAfter >= player0ResourcesBefore {
		t.Fatalf("Expected robber to block resources. Before: %d, After: %d", 
			player0ResourcesBefore, player0ResourcesAfter)
	}
}

func TestRobberDoesNotBlockOtherTiles(t *testing.T) {
	b := NewDesertBoard()
	
	// Set up two specific tiles with different dice numbers
	wheatTileCoord := TileCoord{X: 2, Y: 3}
	b.tiles[wheatTileCoord] = Tile{Terrain: TerrainWheat, DiceNumber: 6}
	
	oreTileCoord := TileCoord{X: 1, Y: 4}
	b.tiles[oreTileCoord] = Tile{Terrain: TerrainOre, DiceNumber: 8}
	
	// Place a settlement at a cross adjacent to both tiles
	cross := CrossCoord{X: 2, Y: 4}
	b.settlements[cross] = 0 // player 0
	
	// Place robber on wheat tile
	b.PlaceRobber(wheatTileCoord)
	
	// Generate resources for dice 6 (wheat tile - should be blocked)
	resources1 := b.GenerateResources(6)
	if len(resources1[0]) > 0 {
		t.Fatal("Robber should block resources from wheat tile")
	}
	
	// Generate resources for dice 8 (ore tile - should NOT be blocked)
	resources2 := b.GenerateResources(8)
	if len(resources2[0]) == 0 {
		t.Fatal("Robber should NOT block resources from ore tile")
	}
}