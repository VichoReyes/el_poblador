package game

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

type phaseWithOptions struct {
	game       *Game
	playerTurn int
	options    []string
	selected   int
}

func (p *phaseWithOptions) BoardCursor() interface{} {
	return nil
}

func (p *phaseWithOptions) Menu() string {
	paddedOptions := make([]string, len(p.options))
	playerColor := p.game.players[p.playerTurn].color
	for i, option := range p.options {
		if i == p.selected {
			paddedOptions[i] = fmt.Sprintf("\033[38;5;%dm>\033[0m %s", playerColor, option)
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

func PhaseDiceRoll(game *Game, playerTurn int) Phase {
	return &phaseDiceRoll{
		phaseWithOptions: phaseWithOptions{
			game:       game,
			playerTurn: playerTurn,
			options:    []string{"Roll", "Play Knight"},
		},
	}
}

func (p *phaseDiceRoll) Confirm() Phase {
	switch p.selected {
	case 0:
		return rollDice(p.game, p.playerTurn)
	case 1:
		panic("Play Knight not implemented")
	default:
		panic("Invalid option selected")
	}
}

func rollDice(game *Game, playerTurn int) Phase {
	game.lastDice = [2]int{rand.IntN(6) + 1, rand.IntN(6) + 1}
	sum := game.lastDice[0] + game.lastDice[1]
	if sum == 7 {
		// TODO: discarding of > 7 cards
		// also TODO: implement robber
		// for now go to idle phase
		return PhaseIdle(game, playerTurn)
	}
	generatedResources := game.board.GenerateResources(sum)
	for player, resources := range generatedResources {
		for _, r := range resources {
			game.players[player].AddResource(r)
		}
	}
	return PhaseIdle(game, playerTurn)
}

func (p *phaseDiceRoll) Cancel() Phase {
	return p
}

func (p *phaseDiceRoll) HelpText() string {
	return fmt.Sprintf("%s's turn. Time to roll the dice",
		p.game.players[p.playerTurn].Name)
}

type phaseIdle struct {
	phaseWithOptions
}

func PhaseIdle(game *Game, playerTurn int) Phase {
	return &phaseIdle{
		phaseWithOptions: phaseWithOptions{
			game:       game,
			playerTurn: playerTurn,
			options:    []string{"Build", "Trade", "Play Development Card", "End Turn"},
		},
	}
}

func (p *phaseIdle) Confirm() Phase {
	return p
}

func (p *phaseIdle) Cancel() Phase {
	return p
}

func (p *phaseIdle) HelpText() string {
	return fmt.Sprintf("%s's turn. What do you want to do?", p.game.players[p.playerTurn].Name)
}
