package board

import "slices"

// Board represents the game board
type Board struct {
	tiles map[TileCoord]Tile
	// roads and settlements are indexed by player id
	roads        map[PathCoord]int
	settlements  map[CrossCoord]int
	cityUpgrades map[CrossCoord]int // tracks which settlements have been upgraded to cities
	playerRender func(int, string) string
	robber       TileCoord
}

// GetRobber returns the current robber position
func (b *Board) GetRobber() TileCoord {
	return b.robber
}

// SetRobber sets the robber to the given tile coordinate
// also returns the player ids of the players that can be stolen from
func (b *Board) PlaceRobber(coord TileCoord) []int {
	b.robber = coord
	playerIds := make([]int, 0)
	for settlement, playerId := range b.settlements {
		if slices.Contains(settlement.adjacentTileCoords(), coord) {
			playerIds = append(playerIds, playerId)
		}
	}
	return playerIds
}

func (b *Board) ValidCrossCoord() CrossCoord {
	if coord, ok := NewCrossCoord(2, 4); ok {
		return coord
	}
	panic("valid cross coord wrong")
}

func (b *Board) CanPlaceSettlement(coord CrossCoord) bool {
	if _, ok := b.settlements[coord]; ok {
		return false
	}
	for _, neighbor := range coord.Neighbors() {
		if _, ok := b.settlements[neighbor]; ok {
			return false
		}
	}
	return true
}

func (b *Board) AdjacentTiles(coord CrossCoord) []Tile {
	tileCoords := coord.adjacentTileCoords()
	tiles := make([]Tile, len(tileCoords))
	for i, tileCoord := range tileCoords {
		tiles[i] = b.tiles[tileCoord]
	}
	return tiles
}

func (b *Board) SetSettlement(coord CrossCoord, playerId int) bool {
	if !b.CanPlaceSettlement(coord) {
		return false
	}
	b.settlements[coord] = playerId
	return true
}

// GenerateResources generates the resources for a given sum
// returns a map of player id to resources
func (b *Board) GenerateResources(sum int) map[int][]ResourceType {
	resources := make(map[int][]ResourceType)

	// First, generate resources for settlements (1 resource each)
	for crossCoord, playerId := range b.settlements {
		tileCoords := crossCoord.adjacentTileCoords()
		for _, tileCoord := range tileCoords {
			// Skip resource generation if robber is on this tile
			if tileCoord == b.robber {
				continue
			}
			tile := b.tiles[tileCoord]
			if tile.DiceNumber == sum {
				resource, ok := TileResource(tile)
				if ok {
					resources[playerId] = append(resources[playerId], resource)
				}
			}
		}
	}

	// Then, generate additional resources for cities (1 more resource each)
	for crossCoord, playerId := range b.cityUpgrades {
		tileCoords := crossCoord.adjacentTileCoords()
		for _, tileCoord := range tileCoords {
			// Skip resource generation if robber is on this tile
			if tileCoord == b.robber {
				continue
			}
			tile := b.tiles[tileCoord]
			if tile.DiceNumber == sum {
				resource, ok := TileResource(tile)
				if ok {
					resources[playerId] = append(resources[playerId], resource)
				}
			}
		}
	}

	return resources
}

func (b *Board) SetRoad(coord PathCoord, playerId int) {
	b.roads[coord] = playerId
}

// HasRoadConnected checks if a player has a road connected to a specific crossing
func (b *Board) HasRoadConnected(cross CrossCoord, playerId int) bool {
	neighbors := cross.Neighbors()
	for _, neighbor := range neighbors {
		pathCoord := NewPathCoord(cross, neighbor)
		if roadPlayerId, exists := b.roads[pathCoord]; exists && roadPlayerId == playerId {
			return true
		}
	}
	return false
}

// CanPlaceRoad checks if a road can be placed at the given path coordinate
func (b *Board) CanPlaceRoad(coord PathCoord, playerId int) bool {
	// Check if road already exists
	if _, ok := b.roads[coord]; ok {
		return false
	}

	// Check if player has a settlement/city at one of the endpoints
	if settlementPlayerId, ok := b.settlements[coord.From]; ok && settlementPlayerId == playerId {
		return true
	}
	if settlementPlayerId, ok := b.settlements[coord.To]; ok && settlementPlayerId == playerId {
		return true
	}

	// Check if player has a road connected to one of the endpoints
	if b.HasRoadConnected(coord.From, playerId) {
		return true
	}
	if b.HasRoadConnected(coord.To, playerId) {
		return true
	}

	return false
}

// HasSettlementAt checks if a player has a settlement at a specific crossing
func (b *Board) HasSettlementAt(cross CrossCoord, playerId int) bool {
	if settlementPlayerId, ok := b.settlements[cross]; ok {
		return settlementPlayerId == playerId
	}
	return false
}

// CanPlaceSettlementForPlayer checks if a player can place a settlement at a specific crossing
func (b *Board) CanPlaceSettlementForPlayer(cross CrossCoord, playerId int) bool {
	// First check if the basic placement rules are satisfied
	if !b.CanPlaceSettlement(cross) {
		return false
	}

	// Then check if the player has a road connected to this crossing
	return b.HasRoadConnected(cross, playerId)
}

// CanUpgradeToCity checks if a settlement can be upgraded to a city
func (b *Board) CanUpgradeToCity(coord CrossCoord, playerId int) bool {
	// Check if there's a settlement owned by this player
	if settlementPlayerId, ok := b.settlements[coord]; !ok || settlementPlayerId != playerId {
		return false
	}

	// Check if it's already been upgraded to a city
	if _, ok := b.cityUpgrades[coord]; ok {
		return false
	}

	return true
}

// UpgradeToCity upgrades a settlement to a city
func (b *Board) UpgradeToCity(coord CrossCoord, playerId int) bool {
	if !b.CanUpgradeToCity(coord, playerId) {
		return false
	}
	b.cityUpgrades[coord] = playerId
	return true
}

// CountSettlements counts the number of settlements owned by a player
func (b *Board) CountSettlements(playerId int) int {
	count := 0
	for _, owner := range b.settlements {
		if owner == playerId {
			count++
		}
	}
	return count
}

// CountCities counts the number of cities owned by a player
func (b *Board) CountCities(playerId int) int {
	count := 0
	for _, owner := range b.cityUpgrades {
		if owner == playerId {
			count++
		}
	}
	return count
}
