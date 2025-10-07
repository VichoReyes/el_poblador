package game

import (
	"bytes"
	"el_poblador/board"
	"encoding/gob"
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"time"
)

type phaseDiceRoll struct {
	phaseWithOptions
	invalid string
}

func PhaseDiceRoll(game *Game) Phase {
	return &phaseDiceRoll{
		phaseWithOptions: phaseWithOptions{
			game:    game,
			options: []string{"Roll", "Play Knight", "Save & Quit"},
		},
	}
}

func (p *phaseDiceRoll) Confirm() Phase {
	switch p.selected {
	case 0:
		return rollDice(p.game)
	case 1:
		// Play Knight card
		player := &p.game.Players[p.game.PlayerTurn]
		if player.PlayDevCard(DevCardKnight) {
			return PhasePlaceRobber(p.game, p)
		}
		p.invalid = "You don't have a Knight card"
		return p
	case 2:
		// Save & Quit
		if err := saveGameState(p.game); err != nil {
			p.invalid = fmt.Sprintf("Save failed: %v", err)
			return p
		}
		p.game.shouldQuit = true
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
	game.LastDice = [2]int{rand.IntN(6) + 1, rand.IntN(6) + 1}
	sum := game.LastDice[0] + game.LastDice[1]
	if sum == 7 {
		// TODO: discarding of > 7 cards
		return PhasePlaceRobber(game, PhaseIdle(game))
	}
	generatedResources := game.Board.GenerateResources(sum)
	for player, resources := range generatedResources {
		for _, r := range resources {
			game.Players[player].AddResource(r)
		}
	}

	// Log resource generation for each player if they received any
	for playerId, resources := range generatedResources {
		if len(resources) > 0 {
			// Count resources by type
			resourceCounts := make(map[board.ResourceType]int)
			for _, resource := range resources {
				resourceCounts[resource]++
			}

			// Build resource description
			var resourceParts []string
			for resourceType, count := range resourceCounts {
				resourceParts = append(resourceParts, fmt.Sprintf("%d %s", count, resourceType))
			}

			game.LogAction(fmt.Sprintf("%s gained %s from dice roll (%d)",
				game.Players[playerId].RenderName(),
				strings.Join(resourceParts, ", "),
				sum))
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
		p.game.PlayerTurn++
		p.game.PlayerTurn %= len(p.game.Players)
		nextPlayer := &p.game.Players[p.game.PlayerTurn]
		p.game.LogAction(fmt.Sprintf("Turn passed to %s", nextPlayer.RenderName()))
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

func saveGameState(g *Game) error {
	filename := fmt.Sprintf("game_save_%s.gob", time.Now().Format("2006-01-02_15-04-05"))
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(g); err != nil {
		return fmt.Errorf("encoding failed: %w", err)
	}
	if err := os.WriteFile(filename, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("write failed: %w", err)
	}
	return nil
}
