package board

import (
	"testing"
)

func TestResourceGenerationWithSettlementsAndCities(t *testing.T) {
	// Create a desert board and modify specific tiles for testing
	board := NewDesertBoard()

	// Set up a specific tile with wheat and dice number 6
	wheatTileCoord := TileCoord{X: 2, Y: 3}
	board.tiles[wheatTileCoord] = Tile{Terrain: TerrainWheat, DiceNumber: 6}

	// Set up a specific tile with ore and dice number 6
	oreTileCoord := TileCoord{X: 1, Y: 4}
	board.tiles[oreTileCoord] = Tile{Terrain: TerrainOre, DiceNumber: 6}

	// Place a settlement at a crossing adjacent to both tiles
	settlementCoord := CrossCoord{X: 2, Y: 4}
	board.settlements[settlementCoord] = 0 // player 0

	// Test resource generation for dice roll 6
	resources := board.GenerateResources(6)

	// Should get 1 wheat and 1 ore from the settlement
	if len(resources[0]) != 2 {
		t.Errorf("Expected 2 resources for player 0, got %d", len(resources[0]))
	}

	// Check that we have both wheat and ore
	wheatCount := 0
	oreCount := 0
	for _, resource := range resources[0] {
		if resource == ResourceWheat {
			wheatCount++
		} else if resource == ResourceOre {
			oreCount++
		}
	}

	if wheatCount != 1 {
		t.Errorf("Expected 1 wheat, got %d", wheatCount)
	}
	if oreCount != 1 {
		t.Errorf("Expected 1 ore, got %d", oreCount)
	}

	// Now upgrade the settlement to a city
	board.cityUpgrades[settlementCoord] = 0

	// Test resource generation again for dice roll 6
	resources = board.GenerateResources(6)

	// Should get 2 wheat and 2 ore from the city (1 from settlement + 1 from city upgrade)
	if len(resources[0]) != 4 {
		t.Errorf("Expected 4 resources for player 0 with city, got %d", len(resources[0]))
	}

	// Check that we have both wheat and ore doubled
	wheatCount = 0
	oreCount = 0
	for _, resource := range resources[0] {
		if resource == ResourceWheat {
			wheatCount++
		} else if resource == ResourceOre {
			oreCount++
		}
	}

	if wheatCount != 2 {
		t.Errorf("Expected 2 wheat with city, got %d", wheatCount)
	}
	if oreCount != 2 {
		t.Errorf("Expected 2 ore with city, got %d", oreCount)
	}
}

func TestCityUpgradeFunctionality(t *testing.T) {
	board := NewDesertBoard()

	// Place a settlement
	settlementCoord := CrossCoord{X: 2, Y: 4}
	board.settlements[settlementCoord] = 0 // player 0

	// Test that we can upgrade the settlement to a city
	if !board.CanUpgradeToCity(settlementCoord, 0) {
		t.Error("Should be able to upgrade own settlement to city")
	}

	// Test that another player cannot upgrade it
	if board.CanUpgradeToCity(settlementCoord, 1) {
		t.Error("Should not be able to upgrade another player's settlement")
	}

	// Test that we cannot upgrade a non-existent settlement
	emptyCoord := CrossCoord{X: 3, Y: 5}
	if board.CanUpgradeToCity(emptyCoord, 0) {
		t.Error("Should not be able to upgrade non-existent settlement")
	}

	// Perform the upgrade
	if !board.UpgradeToCity(settlementCoord, 0) {
		t.Error("Should be able to upgrade settlement to city")
	}

	// Test that we cannot upgrade the same settlement again
	if board.CanUpgradeToCity(settlementCoord, 0) {
		t.Error("Should not be able to upgrade a settlement that's already a city")
	}

	// Test that the upgrade operation fails if we try again
	if board.UpgradeToCity(settlementCoord, 0) {
		t.Error("Should not be able to upgrade a settlement that's already a city")
	}
}

func TestRoadPlacementValidation(t *testing.T) {
	board := NewDesertBoard()

	// Place a settlement
	settlementCoord := CrossCoord{X: 2, Y: 4}
	board.settlements[settlementCoord] = 0 // player 0

	// Place a road connected to the settlement
	neighborCoord := CrossCoord{X: 2, Y: 3}
	roadCoord := NewPathCoord(settlementCoord, neighborCoord)
	board.roads[roadCoord] = 0

	// Test that we can place a road connected to our existing road
	nextNeighborCoord := CrossCoord{X: 2, Y: 2}
	newRoadCoord := NewPathCoord(neighborCoord, nextNeighborCoord)

	if !board.CanPlaceRoad(newRoadCoord, 0) {
		t.Error("Should be able to place road connected to existing road")
	}

	// Test that we cannot place a road if it already exists
	if board.CanPlaceRoad(roadCoord, 0) {
		t.Error("Should not be able to place road where one already exists")
	}

	// Test that we can place a road connected to our settlement
	anotherNeighborCoord := CrossCoord{X: 3, Y: 4}
	settlementRoadCoord := NewPathCoord(settlementCoord, anotherNeighborCoord)

	if !board.CanPlaceRoad(settlementRoadCoord, 0) {
		t.Error("Should be able to place road connected to own settlement")
	}

	// Test that another player cannot place a road connected to our settlement
	if board.CanPlaceRoad(settlementRoadCoord, 1) {
		t.Error("Should not be able to place road connected to another player's settlement")
	}

	// Test that we cannot place a road disconnected from our network
	isolatedCoord1 := CrossCoord{X: 4, Y: 6}
	isolatedCoord2 := CrossCoord{X: 4, Y: 7}
	isolatedRoadCoord := NewPathCoord(isolatedCoord1, isolatedCoord2)

	if board.CanPlaceRoad(isolatedRoadCoord, 0) {
		t.Error("Should not be able to place road disconnected from network")
	}
}

func TestSettlementPlacementValidation(t *testing.T) {
	board := NewDesertBoard()

	// Place a settlement
	settlementCoord := CrossCoord{X: 2, Y: 4}
	board.settlements[settlementCoord] = 0 // player 0

	// Place a road connected to the settlement
	neighborCoord := CrossCoord{X: 2, Y: 3}
	roadCoord := NewPathCoord(settlementCoord, neighborCoord)
	board.roads[roadCoord] = 0

	// Place a road to the location where we want to place the settlement
	roadEndCoord := CrossCoord{X: 2, Y: 2}
	roadToSettlement := NewPathCoord(neighborCoord, roadEndCoord)
	board.roads[roadToSettlement] = 0

	// Test that we can place a settlement connected to our road
	if !board.CanPlaceSettlementForPlayer(roadEndCoord, 0) {
		t.Error("Should be able to place settlement connected to own road")
	}

	// Test that we cannot place a settlement at the same location
	if board.CanPlaceSettlementForPlayer(settlementCoord, 0) {
		t.Error("Should not be able to place settlement where one already exists")
	}

	// Test that we cannot place a settlement adjacent to existing settlement
	adjacentCoord := CrossCoord{X: 2, Y: 5}
	if board.CanPlaceSettlementForPlayer(adjacentCoord, 0) {
		t.Error("Should not be able to place settlement adjacent to existing settlement")
	}

	// Test that we cannot place a settlement disconnected from our road network
	isolatedCoord := CrossCoord{X: 4, Y: 6}
	if board.CanPlaceSettlementForPlayer(isolatedCoord, 0) {
		t.Error("Should not be able to place settlement disconnected from road network")
	}

	// Test that another player cannot place a settlement connected to our road
	if board.CanPlaceSettlementForPlayer(roadEndCoord, 1) {
		t.Error("Should not be able to place settlement connected to another player's road")
	}
}

func TestRoadConnectionDetection(t *testing.T) {
	board := NewDesertBoard()

	// Place a settlement
	settlementCoord := CrossCoord{X: 2, Y: 4}
	board.settlements[settlementCoord] = 0 // player 0

	// Test that settlement has no roads connected initially
	if board.HasRoadConnected(settlementCoord, 0) {
		t.Error("Settlement should not have roads connected initially")
	}

	// Place a road connected to the settlement
	neighborCoord := CrossCoord{X: 2, Y: 3}
	roadCoord := NewPathCoord(settlementCoord, neighborCoord)
	board.roads[roadCoord] = 0

	// Test that settlement now has a road connected
	if !board.HasRoadConnected(settlementCoord, 0) {
		t.Error("Settlement should have road connected after placing road")
	}

	// Test that another player doesn't have roads connected
	if board.HasRoadConnected(settlementCoord, 1) {
		t.Error("Another player should not have roads connected to this settlement")
	}

	// Test that neighbor crossing has road connected
	if !board.HasRoadConnected(neighborCoord, 0) {
		t.Error("Neighbor crossing should have road connected")
	}

	// Test that a distant crossing has no roads connected
	distantCoord := CrossCoord{X: 4, Y: 6}
	if board.HasRoadConnected(distantCoord, 0) {
		t.Error("Distant crossing should not have roads connected")
	}
}

func TestSettlementOwnershipDetection(t *testing.T) {
	board := NewDesertBoard()

	// Place a settlement
	settlementCoord := CrossCoord{X: 2, Y: 4}
	board.settlements[settlementCoord] = 0 // player 0

	// Test that player 0 owns the settlement
	if !board.HasSettlementAt(settlementCoord, 0) {
		t.Error("Player 0 should own the settlement")
	}

	// Test that player 1 doesn't own the settlement
	if board.HasSettlementAt(settlementCoord, 1) {
		t.Error("Player 1 should not own the settlement")
	}

	// Test that no player owns an empty location
	emptyCoord := CrossCoord{X: 3, Y: 5}
	if board.HasSettlementAt(emptyCoord, 0) {
		t.Error("No player should own an empty location")
	}
}
