package game

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type phaseWithOptions struct {
	game     *Game
	options  []string
	selected int
}

func (p *phaseWithOptions) BoardCursor() interface{} {
	return nil
}

func (p *phaseWithOptions) Menu() string {
	paddedOptions := make([]string, len(p.options))
	player := p.game.players[p.game.playerTurn]
	for i, option := range p.options {
		if i == p.selected {
			paddedOptions[i] = player.Render("> ") + option
		} else {
			paddedOptions[i] = fmt.Sprintf(" %s", option)
		}
	}
	return strings.Join(paddedOptions, "\n")
}

func (p *phaseWithOptions) MoveCursor(direction string) {
	switch direction {
	case "up":
		p.selected--
	case "down":
		p.selected++
	}
	p.selected = (p.selected + len(p.options)) % len(p.options)
}

type phaseDiceRoll struct {
	phaseWithOptions
}

func PhaseDiceRoll(game *Game) Phase {
	return &phaseDiceRoll{
		phaseWithOptions: phaseWithOptions{
			game:    game,
			options: []string{"Roll", "Play Knight"},
		},
	}
}

func (p *phaseDiceRoll) Confirm() Phase {
	switch p.selected {
	case 0:
		return rollDice(p.game)
	case 1:
		panic("Play Knight not implemented")
	default:
		panic("Invalid option selected")
	}
}

func rollDice(game *Game) Phase {
	game.lastDice = [2]int{rand.IntN(6) + 1, rand.IntN(6) + 1}
	sum := game.lastDice[0] + game.lastDice[1]
	if sum == 7 {
		// TODO: discarding of > 7 cards
		// also TODO: implement robber
		// for now go to idle phase
		return PhaseIdle(game)
	}
	generatedResources := game.board.GenerateResources(sum)
	for player, resources := range generatedResources {
		for _, r := range resources {
			game.players[player].AddResource(r)
		}
	}
	return PhaseIdle(game)
}

func (p *phaseDiceRoll) HelpText() string {
	return "Time to roll the dice"
}

type phaseIdle struct {
	phaseWithOptions
}

func PhaseIdle(game *Game) Phase {
	return &phaseIdle{
		phaseWithOptions: phaseWithOptions{
			game:    game,
			options: []string{"Build", "Trade", "Play Development Card", "End Turn"},
		},
	}
}

func (p *phaseIdle) Confirm() Phase {
	switch p.selected {
	case 0: // Build
		return PhaseBuilding(p.game, p)
	case 1: // Trade
		panic("Trade not implemented")
	case 2: // Play Development Card
		panic("Play Development Card not implemented")
	case 3: // End Turn
		p.game.playerTurn++
		p.game.playerTurn %= len(p.game.players)
		return PhaseIdle(p.game)
	default:
		panic("Invalid option selected")
	}
}

func (p *phaseIdle) HelpText() string {
	return "What do you want to do?"
}

type phaseBuilding struct {
	phaseWithOptions
	previousPhase Phase
}

func PhaseBuilding(game *Game, previousPhase Phase) Phase {
	player := game.players[game.playerTurn]

	// Build the list of available building options
	var options []string
	strikethrough := lipgloss.NewStyle().Strikethrough(true)

	if player.CanBuildRoad() {
		options = append(options, "Road")
	} else {
		options = append(options, strikethrough.Render("Road"))
	}

	if player.CanBuildSettlement() {
		options = append(options, "Settlement")
	} else {
		options = append(options, strikethrough.Render("Settlement"))
	}

	if player.CanBuildCity() {
		options = append(options, "City")
	} else {
		options = append(options, strikethrough.Render("City"))
	}

	if player.CanBuyDevelopmentCard() {
		options = append(options, "Development Card")
	} else {
		options = append(options, strikethrough.Render("Development Card"))
	}

	options = append(options, "Cancel (or 'esc')")

	return &phaseBuilding{
		phaseWithOptions: phaseWithOptions{
			game:    game,
			options: options,
		},
		previousPhase: previousPhase,
	}
}

func (p *phaseBuilding) Confirm() Phase {
	player := p.game.players[p.game.playerTurn]

	switch p.selected {
	case 0: // Road
		if player.CanBuildRoad() {
			// TODO: Implement road placement phase
			panic("Road placement not implemented")
		}
		return p
	case 1: // Settlement
		if player.CanBuildSettlement() {
			// TODO: Implement settlement placement phase
			panic("Settlement placement not implemented")
		}
		return p
	case 2: // City
		if player.CanBuildCity() {
			// TODO: Implement city placement phase
			panic("City placement not implemented")
		}
		return p
	case 3: // Development Card
		if player.CanBuyDevelopmentCard() {
			// TODO: Implement development card purchase
			panic("Development card purchase not implemented")
		}
		return p
	case 4: // Cancel
		return p.previousPhase
	default:
		panic("Invalid option selected")
	}
}

func (p *phaseBuilding) Cancel() Phase {
	return p.previousPhase
}

func (p *phaseBuilding) HelpText() string {
	return "Choose what to build"
}
