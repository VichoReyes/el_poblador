package game

import (
	"math/rand/v2"
)

type phaseDiceRoll struct {
	phaseWithOptions
	invalid string
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
		// Play Knight card
		player := &p.game.players[p.game.playerTurn]
		if player.PlayDevCard(DevCardKnight) {
			return PhasePlaceRobber(p.game, p)
		}
		p.invalid = "You don't have a Knight card"
		return p
	default:
		panic("Invalid option selected")
	}
}

func (p *phaseDiceRoll) HelpText() string {
	if p.invalid != "" {
		return p.invalid
	}
	return "Time to roll the dice"
}

func rollDice(game *Game) Phase {
	game.lastDice = [2]int{rand.IntN(6) + 1, rand.IntN(6) + 1}
	sum := game.lastDice[0] + game.lastDice[1]
	if sum == 7 {
		// TODO: discarding of > 7 cards
		return PhasePlaceRobber(game, PhaseIdle(game))
	}
	generatedResources := game.board.GenerateResources(sum)
	for player, resources := range generatedResources {
		for _, r := range resources {
			game.players[player].AddResource(r)
		}
	}
	return PhaseIdle(game)
}

type phaseIdle struct {
	phaseWithOptions
	notification string
}

func PhaseIdle(game *Game) Phase {
	return PhaseIdleWithNotification(game, "")
}

func PhaseIdleWithNotification(game *Game, notification string) Phase {
	return &phaseIdle{
		phaseWithOptions: phaseWithOptions{
			game:    game,
			options: []string{"Build", "Trade", "Play Development Card", "End Turn"},
		},
		notification: notification,
	}
}

func (p *phaseIdle) Confirm() Phase {
	switch p.selected {
	case 0: // Build
		return PhaseBuilding(p.game, p)
	case 1: // Trade
		p.notification = "Trade not yet implemented."
		return p
	case 2: // Play Development Card
		return PhasePlayDevelopmentCard(p.game, p)
	case 3: // End Turn
		p.game.playerTurn++
		p.game.playerTurn %= len(p.game.players)
		return PhaseDiceRoll(p.game)
	default:
		panic("Invalid option selected")
	}
}

func (p *phaseIdle) HelpText() string {
	if p.notification != "" {
		return p.notification + " Anything else?"
	}
	return "What do you want to do?"
}
