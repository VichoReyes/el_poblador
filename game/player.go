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

func (p *Player) Render(s string) string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(fmt.Sprintf("%d", p.color)))
	return style.Render(s)
}
