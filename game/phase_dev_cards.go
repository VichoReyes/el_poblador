package game

import (
	"el_poblador/board"
	"fmt"
)

type phasePlayDevelopmentCard struct {
	phaseWithOptions
	previousPhase Phase
}

func PhasePlayDevelopmentCard(game *Game, previousPhase Phase) Phase {
	player := game.players[game.playerTurn]

	// Build options based on available development cards
	var options []string

	for _, card := range player.hiddenDevCards {
		options = append(options, card.String())
	}

	options = append(options, "Cancel")

	return &phasePlayDevelopmentCard{
		phaseWithOptions: phaseWithOptions{
			game:    game,
			options: options,
		},
		previousPhase: previousPhase,
	}
}

func (p *phasePlayDevelopmentCard) Confirm() Phase {
	player := &p.game.players[p.game.playerTurn]
	numCards := len(player.hiddenDevCards)
	if p.selected == numCards {
		return p.previousPhase
	}
	card := player.hiddenDevCards[p.selected]

	switch card {
	case DevCardKnight:
		player.PlayDevCard(card)
		return PhasePlaceRobber(p.game, PhaseIdle(p.game))
	case DevCardRoadBuilding:
		player.PlayDevCard(card)
		return PhaseRoadBuilding(p.game)
	case DevCardMonopoly:
		player.PlayDevCard(card)
		return PhaseMonopoly(p.game, PhaseIdle(p.game))
	case DevCardYearOfPlenty:
		player.PlayDevCard(card)
		return PhaseYearOfPlenty(p.game, PhaseIdle(p.game))
	case DevCardVictoryPoint:
		return p.previousPhase
	default:
		panic("This card does not exist")
	}
}

func (p *phasePlayDevelopmentCard) Cancel() Phase {
	return p.previousPhase
}

func (p *phasePlayDevelopmentCard) HelpText() string {
	return "Choose a development card to play"
}

type phaseMonopoly struct {
	phaseWithOptions
	previousPhase Phase
}

func PhaseMonopoly(game *Game, previousPhase Phase) Phase {
	resourceOptions := make([]string, len(board.RESOURCE_TYPES))
	for i, resourceType := range board.RESOURCE_TYPES {
		resourceOptions[i] = string(resourceType)
	}
	resourceOptions = append(resourceOptions, "Cancel")

	return &phaseMonopoly{
		phaseWithOptions: phaseWithOptions{
			game:    game,
			options: resourceOptions,
		},
		previousPhase: previousPhase,
	}
}

func (p *phaseMonopoly) Confirm() Phase {
	if p.selected == len(board.RESOURCE_TYPES) {
		return p.previousPhase
	}

	selectedResource := board.RESOURCE_TYPES[p.selected]
	currentPlayer := p.game.players[p.game.playerTurn]

	totalCollected := 0
	for i, player := range p.game.players {
		if i != p.game.playerTurn {
			count := player.resources[selectedResource]
			if count > 0 {
				player.resources[selectedResource] = 0
				totalCollected += count
			}
		}
	}

	if totalCollected > 0 {
		currentPlayer.resources[selectedResource] += totalCollected
		return PhaseIdleWithNotification(p.game, fmt.Sprintf("Collected %d %s from other players!", totalCollected, selectedResource))
	} else {
		return PhaseIdleWithNotification(p.game, "No resources collected - nobody had any!")
	}
}

func (p *phaseMonopoly) Cancel() Phase {
	return p.previousPhase
}

func (p *phaseMonopoly) HelpText() string {
	return "Select a resource type to collect from all players"
}

type phaseYearOfPlenty struct {
	phaseWithOptions
	previousPhase     Phase
	selectedCount     int
	selectedResources [2]board.ResourceType
}

func PhaseYearOfPlenty(game *Game, previousPhase Phase) Phase {
	resourceOptions := make([]string, len(board.RESOURCE_TYPES))
	for i, resourceType := range board.RESOURCE_TYPES {
		resourceOptions[i] = string(resourceType)
	}
	resourceOptions = append(resourceOptions, "Cancel")

	return &phaseYearOfPlenty{
		phaseWithOptions: phaseWithOptions{
			game:    game,
			options: resourceOptions,
		},
		previousPhase: previousPhase,
		selectedCount: 0,
	}
}

func (p *phaseYearOfPlenty) Confirm() Phase {
	if p.selected == len(board.RESOURCE_TYPES) {
		return p.previousPhase
	}

	selectedResource := board.RESOURCE_TYPES[p.selected]
	p.selectedResources[p.selectedCount] = selectedResource
	p.selectedCount++

	if p.selectedCount < 2 {
		// Still need to select more resources
		return p
	}

	// Both resources selected, give them to the player
	currentPlayer := &p.game.players[p.game.playerTurn]
	for _, resource := range p.selectedResources {
		currentPlayer.AddResource(resource)
	}

	return PhaseIdleWithNotification(p.game, fmt.Sprintf("Gained %s and %s from the bank!", p.selectedResources[0], p.selectedResources[1]))
}

func (p *phaseYearOfPlenty) Cancel() Phase {
	return p.previousPhase
}

func (p *phaseYearOfPlenty) HelpText() string {
	if p.selectedCount == 0 {
		return "Select first resource to gain from the bank"
	}
	return fmt.Sprintf("Selected: %s. Select second resource to gain from the bank", p.selectedResources[0])
}