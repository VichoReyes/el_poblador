package game

import (
	"el_poblador/board"
	"fmt"
	"strings"
)

// Trading System
//
// This file implements resource trading in El Poblador. The trading system uses a two-phase
// approach where players first specify what they want to offer, then what they want to receive.
// Validation only happens after the complete trade proposal is known, making the system extensible
// for different trade types.
//
// Phase Flow:
//   phaseIdle → phaseTradeOffer → phaseTradeSelectReceive → validateAndExecuteTrade → phaseIdle
//
// UI Controls:
//   - Up/Down: Navigate between resources
//   - Left/Right: Adjust resource amounts
//   - Enter: Confirm and proceed to next phase
//   - Esc: Cancel and return to previous phase
//
// Trade Types:
//   - Bank Trade (4:1): Offer exactly 4 of one resource for exactly 1 of another
//   - Harbor Trade (future): Offer 3:1 or 2:1 depending on harbor access
//   - Player Trade (future): Any combination of resources offered/requested
//
// Adding New Trade Types:
//   1. Add detection function (e.g., isHarborTrade) similar to isBankTrade
//   2. Add execution logic in validateAndExecuteTrade
//   3. For player trades, add new phases for partner selection and negotiation
//
// Design Notes:
//   - phaseTradeOffer has no previousPhase (always returns to phaseIdle)
//   - phaseTradeSelectReceive keeps reference to offer phase for cancel preservation
//   - Validation is deferred until complete trade is known (offer + request)
//   - Trade type detection happens in isBankTrade, isHarborTrade (future), etc.
//   - Both phases implement PhaseCancelable for Esc key handling

type phaseTradeOffer struct {
	game     *Game
	offer    map[board.ResourceType]int
	selected int
}

func PhaseTradeOffer(game *Game) Phase {
	offer := make(map[board.ResourceType]int)
	for _, resourceType := range board.RESOURCE_TYPES {
		offer[resourceType] = 0
	}

	return &phaseTradeOffer{
		game:     game,
		offer:    offer,
		selected: 0,
	}
}

func (p *phaseTradeOffer) MoveCursor(direction string) {
	numResources := len(board.RESOURCE_TYPES)

	switch direction {
	case "up":
		p.selected--
		if p.selected < 0 {
			p.selected = numResources - 1
		}
	case "down":
		p.selected++
		if p.selected >= numResources {
			p.selected = 0
		}
	case "left":
		resourceType := board.RESOURCE_TYPES[p.selected]
		if p.offer[resourceType] > 0 {
			p.offer[resourceType]--
		}
	case "right":
		resourceType := board.RESOURCE_TYPES[p.selected]
		player := &p.game.Players[p.game.PlayerTurn]
		maxAvailable := player.Resources[resourceType]
		if p.offer[resourceType] < maxAvailable {
			p.offer[resourceType]++
		}
	}
}

func (p *phaseTradeOffer) Confirm() Phase {
	totalOffered := 0
	for _, amount := range p.offer {
		totalOffered += amount
	}
	if totalOffered == 0 {
		return p
	}
	return PhaseTradeSelectReceive(p.game, p.offer, p)
}

func (p *phaseTradeOffer) Cancel() Phase {
	return PhaseIdle(p.game)
}

func (p *phaseTradeOffer) BoardCursor() interface{} {
	return nil
}

func (p *phaseTradeOffer) Menu() string {
	var lines []string
	player := &p.game.Players[p.game.PlayerTurn]

	lines = append(lines, "What do you want to offer?")
	lines = append(lines, "")

	for i, resourceType := range board.RESOURCE_TYPES {
		amount := p.offer[resourceType]
		maxAvailable := player.Resources[resourceType]

		line := fmt.Sprintf("%s:  %d / %d", resourceType, amount, maxAvailable)
		if i == p.selected {
			line = player.Render("> ") + line
		} else {
			line = "  " + line
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func (p *phaseTradeOffer) HelpText() string {
	return "Use ←/→ to adjust, ↑/↓ to move, Enter to confirm, Esc to cancel"
}

type phaseTradeSelectReceive struct {
	game          *Game
	offer         map[board.ResourceType]int
	request       map[board.ResourceType]int
	previousPhase Phase
	selected      int
}

func PhaseTradeSelectReceive(game *Game, offer map[board.ResourceType]int, offerPhase Phase) Phase {
	request := make(map[board.ResourceType]int)
	for _, resourceType := range board.RESOURCE_TYPES {
		request[resourceType] = 0
	}

	return &phaseTradeSelectReceive{
		game:          game,
		offer:         offer,
		request:       request,
		previousPhase: offerPhase,
		selected:      0,
	}
}

func (p *phaseTradeSelectReceive) MoveCursor(direction string) {
	numResources := len(board.RESOURCE_TYPES)

	switch direction {
	case "up":
		p.selected--
		if p.selected < 0 {
			p.selected = numResources - 1
		}
	case "down":
		p.selected++
		if p.selected >= numResources {
			p.selected = 0
		}
	case "left":
		resourceType := board.RESOURCE_TYPES[p.selected]
		if p.request[resourceType] > 0 {
			p.request[resourceType]--
		}
	case "right":
		resourceType := board.RESOURCE_TYPES[p.selected]
		p.request[resourceType]++
	}
}

func (p *phaseTradeSelectReceive) Confirm() Phase {
	totalRequested := 0
	for _, amount := range p.request {
		totalRequested += amount
	}
	if totalRequested == 0 {
		return p
	}

	return p.validateAndExecuteTrade()
}

func (p *phaseTradeSelectReceive) Cancel() Phase {
	return p.previousPhase
}

func (p *phaseTradeSelectReceive) BoardCursor() interface{} {
	return nil
}

func (p *phaseTradeSelectReceive) Menu() string {
	var lines []string
	player := &p.game.Players[p.game.PlayerTurn]

	lines = append(lines, "What do you want to receive?")
	lines = append(lines, "")

	for i, resourceType := range board.RESOURCE_TYPES {
		amount := p.request[resourceType]

		line := fmt.Sprintf("%s:  %d", resourceType, amount)
		if i == p.selected {
			line = player.Render("> ") + line
		} else {
			line = "  " + line
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func (p *phaseTradeSelectReceive) HelpText() string {
	return "Use ←/→ to adjust, ↑/↓ to move, Enter to confirm, Esc to cancel"
}

func (p *phaseTradeSelectReceive) validateAndExecuteTrade() Phase {
	player := &p.game.Players[p.game.PlayerTurn]

	if tradeType, offeredResource, requestedResource := p.isBankTrade(); tradeType == "bank" {
		if player.Resources[offeredResource] < 4 {
			return PhaseIdleWithNotification(p.game, "Not enough resources for bank trade!")
		}

		player.Resources[offeredResource] -= 4
		player.Resources[requestedResource]++

		p.game.LogAction(fmt.Sprintf("%s traded 4 %s for 1 %s with the bank",
			player.RenderName(), offeredResource, requestedResource))

		return PhaseIdleWithNotification(p.game,
			fmt.Sprintf("Traded 4 %s for 1 %s!", offeredResource, requestedResource))
	}

	return PhaseIdleWithNotification(p.game, "Trade type not yet implemented")
}

func (p *phaseTradeSelectReceive) isBankTrade() (string, board.ResourceType, board.ResourceType) {
	totalOffered := 0
	offeredTypes := 0
	var offeredResource board.ResourceType
	for resourceType, amount := range p.offer {
		if amount > 0 {
			offeredTypes++
			offeredResource = resourceType
			totalOffered += amount
		}
	}

	totalRequested := 0
	requestedTypes := 0
	var requestedResource board.ResourceType
	for resourceType, amount := range p.request {
		if amount > 0 {
			requestedTypes++
			requestedResource = resourceType
			totalRequested += amount
		}
	}

	if totalOffered == 4 && offeredTypes == 1 && totalRequested == 1 && requestedTypes == 1 {
		return "bank", offeredResource, requestedResource
	}

	return "unknown", "", ""
}
