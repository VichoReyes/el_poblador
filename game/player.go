package game

import (
	"el_poblador/board"
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type Player struct {
	Name           string
	color          int // 8 bit color code
	resources      map[board.ResourceType]int
	hiddenDevCards []DevCard
	playedDevCards []DevCard
}

func (p *Player) TotalResources() int {
	total := 0
	for _, amount := range p.resources {
		total += amount
	}
	return total
}

func (p *Player) AddResource(t board.ResourceType) {
	p.resources[t] += 1
}

// HasResources checks if the player has the required resources
func (p *Player) HasResources(required map[board.ResourceType]int) bool {
	for resource, amount := range required {
		if p.resources[resource] < amount {
			return false
		}
	}
	return true
}

// ConsumeResources removes the specified resources from the player
func (p *Player) ConsumeResources(required map[board.ResourceType]int) bool {
	if !p.HasResources(required) {
		return false
	}
	for resource, amount := range required {
		p.resources[resource] -= amount
	}
	return true
}

// CanBuildRoad checks if the player can afford to build a road
func (p *Player) CanBuildRoad() bool {
	required := map[board.ResourceType]int{
		board.ResourceWood:  1,
		board.ResourceBrick: 1,
	}
	return p.HasResources(required)
}

// CanBuildSettlement checks if the player can afford to build a settlement
func (p *Player) CanBuildSettlement() bool {
	required := map[board.ResourceType]int{
		board.ResourceWood:  1,
		board.ResourceBrick: 1,
		board.ResourceWheat: 1,
		board.ResourceSheep: 1,
	}
	return p.HasResources(required)
}

// CanBuildCity checks if the player can afford to build a city
func (p *Player) CanBuildCity() bool {
	required := map[board.ResourceType]int{
		board.ResourceWheat: 2,
		board.ResourceOre:   3,
	}
	return p.HasResources(required)
}

// CanBuyDevelopmentCard checks if the player can afford to buy a development card
func (p *Player) CanBuyDevelopmentCard() bool {
	required := map[board.ResourceType]int{
		board.ResourceWheat: 1,
		board.ResourceOre:   1,
		board.ResourceSheep: 1,
	}
	return p.HasResources(required)
}

// BuyDevelopmentCard consumes resources and returns true if successful
func (p *Player) BuyDevelopmentCard() bool {
	required := map[board.ResourceType]int{
		board.ResourceWheat: 1,
		board.ResourceOre:   1,
		board.ResourceSheep: 1,
	}
	return p.ConsumeResources(required)
}

// PlayDevCard moves a card from hidden to played deck
func (p *Player) PlayDevCard(card DevCard) bool {
	for i, hiddenCard := range p.hiddenDevCards {
		if hiddenCard == card {
			// Remove from hidden deck
			p.hiddenDevCards = append(p.hiddenDevCards[:i], p.hiddenDevCards[i+1:]...)
			// Add to played deck
			p.playedDevCards = append(p.playedDevCards, card)
			return true
		}
	}
	return false
}

// TotalDevCards returns the total number of development cards the player has
func (p *Player) TotalDevCards() int {
	return len(p.hiddenDevCards) + len(p.playedDevCards)
}

// BuildRoad consumes resources and builds a road
func (p *Player) BuildRoad() bool {
	required := map[board.ResourceType]int{
		board.ResourceWood:  1,
		board.ResourceBrick: 1,
	}
	return p.ConsumeResources(required)
}

// BuildSettlement consumes resources and builds a settlement
func (p *Player) BuildSettlement() bool {
	required := map[board.ResourceType]int{
		board.ResourceWood:  1,
		board.ResourceBrick: 1,
		board.ResourceWheat: 1,
		board.ResourceSheep: 1,
	}
	return p.ConsumeResources(required)
}

// BuildCity consumes resources and builds a city
func (p *Player) BuildCity() bool {
	required := map[board.ResourceType]int{
		board.ResourceWheat: 2,
		board.ResourceOre:   3,
	}
	return p.ConsumeResources(required)
}

// VictoryPoints calculates the player's current victory points
func (p *Player) VictoryPoints(game *Game) int {
	points := 0
	
	// Count victory point development cards (both hidden and played)
	for _, card := range p.hiddenDevCards {
		if card == DevCardVictoryPoint {
			points++
		}
	}
	for _, card := range p.playedDevCards {
		if card == DevCardVictoryPoint {
			points++
		}
	}
	
	// Count settlements and cities on the board
	playerID := game.getPlayerID(p)
	if playerID != -1 {
		// Settlements are worth 1 point each
		points += game.board.CountSettlements(playerID)
		// Cities replace settlements, so they're worth 2 points total (not additional 2)
		// But in this codebase cities are stored separately from settlements, so count cities as 1 additional point
		points += game.board.CountCities(playerID)
	}
	
	return points
}

func (p *Player) Render(s string) string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(fmt.Sprintf("%d", p.color)))
	return style.Render(s)
}
