package board

import (
	"slices"

	"github.com/charmbracelet/lipgloss"
)

// Board represents the game board
type Board struct {
	Tiles map[TileCoord]Tile
	// Roads and Settlements are indexed by player id
	Roads        map[PathCoord]int
	Settlements  map[CrossCoord]int
	CityUpgrades map[CrossCoord]int             // tracks which settlements have been upgraded to cities
	PlayerColors map[int]lipgloss.AdaptiveColor // player id to adaptive color for rendering
	Robber       TileCoord
}

// GetRobber returns the current robber position
func (b *Board) GetRobber() TileCoord {
	return b.Robber
}

// SetRobber sets the robber to the given tile coordinate
// also returns the player ids of the players that can be stolen from
func (b *Board) PlaceRobber(coord TileCoord) []int {
	b.Robber = coord
	playerIds := make([]int, 0)
	for settlement, playerId := range b.Settlements {
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

func (b *Board) ValidTileCoord() TileCoord {
	if coord, ok := NewTileCoord(1, 2); ok {
		return coord
	}
	panic("valid tile coord wrong")
}

func (b *Board) CanPlaceSettlement(coord CrossCoord) bool {
	if _, ok := b.Settlements[coord]; ok {
		return false
	}
	for _, neighbor := range coord.Neighbors() {
		if _, ok := b.Settlements[neighbor]; ok {
			return false
		}
	}
	return true
}

func (b *Board) AdjacentTiles(coord CrossCoord) []Tile {
	tileCoords := coord.adjacentTileCoords()
	tiles := make([]Tile, len(tileCoords))
	for i, tileCoord := range tileCoords {
		tiles[i] = b.Tiles[tileCoord]
	}
	return tiles
}

func (b *Board) SetSettlement(coord CrossCoord, playerId int) bool {
	if !b.CanPlaceSettlement(coord) {
		return false
	}
	b.Settlements[coord] = playerId
	return true
}

// GenerateResources generates the resources for a given sum
// returns a map of player id to resources
func (b *Board) GenerateResources(sum int) map[int][]ResourceType {
	resources := make(map[int][]ResourceType)

	// First, generate resources for settlements (1 resource each)
	for crossCoord, playerId := range b.Settlements {
		tileCoords := crossCoord.adjacentTileCoords()
		for _, tileCoord := range tileCoords {
			// Skip resource generation if robber is on this tile
			if tileCoord == b.Robber {
				continue
			}
			tile := b.Tiles[tileCoord]
			if tile.DiceNumber == sum {
				resource, ok := TileResource(tile)
				if ok {
					resources[playerId] = append(resources[playerId], resource)
				}
			}
		}
	}

	// Then, generate additional resources for cities (1 more resource each)
	for crossCoord, playerId := range b.CityUpgrades {
		tileCoords := crossCoord.adjacentTileCoords()
		for _, tileCoord := range tileCoords {
			// Skip resource generation if robber is on this tile
			if tileCoord == b.Robber {
				continue
			}
			tile := b.Tiles[tileCoord]
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
	b.Roads[coord] = playerId
}

// HasRoadConnected checks if a player has a road connected to a specific crossing
func (b *Board) HasRoadConnected(cross CrossCoord, playerId int) bool {
	neighbors := cross.Neighbors()
	for _, neighbor := range neighbors {
		pathCoord := NewPathCoord(cross, neighbor)
		if roadPlayerId, exists := b.Roads[pathCoord]; exists && roadPlayerId == playerId {
			return true
		}
	}
	return false
}

// CanPlaceRoad checks if a road can be placed at the given path coordinate
func (b *Board) CanPlaceRoad(coord PathCoord, playerId int) bool {
	// Check if road already exists
	if _, ok := b.Roads[coord]; ok {
		return false
	}

	// Check if player has a settlement/city at one of the endpoints
	if settlementPlayerId, ok := b.Settlements[coord.From]; ok && settlementPlayerId == playerId {
		return true
	}
	if settlementPlayerId, ok := b.Settlements[coord.To]; ok && settlementPlayerId == playerId {
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
	if settlementPlayerId, ok := b.Settlements[cross]; ok {
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
	if settlementPlayerId, ok := b.Settlements[coord]; !ok || settlementPlayerId != playerId {
		return false
	}

	// Check if it's already been upgraded to a city
	if _, ok := b.CityUpgrades[coord]; ok {
		return false
	}

	return true
}

// UpgradeToCity upgrades a settlement to a city
func (b *Board) UpgradeToCity(coord CrossCoord, playerId int) bool {
	if !b.CanUpgradeToCity(coord, playerId) {
		return false
	}
	b.CityUpgrades[coord] = playerId
	return true
}

// CountSettlements counts the number of settlements owned by a player
func (b *Board) CountSettlements(playerId int) int {
	count := 0
	for _, owner := range b.Settlements {
		if owner == playerId {
			count++
		}
	}
	return count
}

// CountCities counts the number of cities owned by a player
func (b *Board) CountCities(playerId int) int {
	count := 0
	for _, owner := range b.CityUpgrades {
		if owner == playerId {
			count++
		}
	}
	return count
}
