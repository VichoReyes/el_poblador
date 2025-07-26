package game

import (
	"fmt"
	"math/rand/v2"
)

type DiceRollPhase struct {
	game       *Game
	playerTurn int
	options    []string
	selected   int
}

func PhaseDiceRoll(game *Game, playerTurn int) Phase {
	return &DiceRollPhase{
		game:       game,
		playerTurn: playerTurn,
		options:    []string{"Roll", "Play Knight"},
	}
}

func (p *DiceRollPhase) Confirm() Phase {
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

func (p *DiceRollPhase) Cancel() Phase {
	return p
}

func (p *DiceRollPhase) HelpText() string {
	return fmt.Sprintf("%s's turn. Time to roll the dice",
		p.game.players[p.playerTurn].Name)
}

func (p *DiceRollPhase) CurrentCursor() interface{} {
	return nil
}

func (p *DiceRollPhase) MoveCursor(direction string) {
	switch direction {
	case "up":
		p.selected--
	case "down":
		p.selected++
	}
	p.selected = (p.selected + len(p.options)) % len(p.options)
}

type phaseIdle struct {
	game       *Game
	playerTurn int
}

func PhaseIdle(game *Game, playerTurn int) Phase {
	return &phaseIdle{game: game, playerTurn: playerTurn}
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

func (p *phaseIdle) CurrentCursor() interface{} {
	return nil
}

func (p *phaseIdle) MoveCursor(direction string) {
	// TODO
}
