package game

import (
	"el_poblador/board"
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type Player struct {
	Name      string
	color     int // 8 bit color code
	resources map[board.ResourceType]int
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

func (p *Player) Render(s string) string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(fmt.Sprintf("%d", p.color)))
	return style.Render(s)
}
