package game

import (
	"el_poblador/board"
	"fmt"
	"strings"
)

// phaseTradeOffer allows player to specify resources to offer
// Use up/down to move cursor, left/right to adjust amounts
type phaseTradeOffer struct {
	game     *Game
	offer    map[board.ResourceType]int
	selected int // 0-4 for resources, 5 for confirm, 6 for cancel
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
	numOptions := numResources + 2 // resources + confirm + cancel

	switch direction {
	case "up":
		p.selected--
		if p.selected < 0 {
			p.selected = numOptions - 1
		}
	case "down":
		p.selected++
		if p.selected >= numOptions {
			p.selected = 0
		}
	case "left":
		// Decrement resource amount if cursor on a resource
		if p.selected < numResources {
			resourceType := board.RESOURCE_TYPES[p.selected]
			if p.offer[resourceType] > 0 {
				p.offer[resourceType]--
			}
		}
	case "right":
		// Increment resource amount if cursor on a resource
		if p.selected < numResources {
			resourceType := board.RESOURCE_TYPES[p.selected]
			player := &p.game.Players[p.game.PlayerTurn]
			maxAvailable := player.Resources[resourceType]
			if p.offer[resourceType] < maxAvailable {
				p.offer[resourceType]++
			}
		}
	}
}

func (p *phaseTradeOffer) Confirm() Phase {
	numResources := len(board.RESOURCE_TYPES)

	// On "Confirm" button - no validation, just proceed
	if p.selected == numResources {
		// Check if offering anything
		totalOffered := 0
		for _, amount := range p.offer {
			totalOffered += amount
		}
		if totalOffered == 0 {
			// Stay in phase if nothing offered
			return p
		}
		return PhaseTradeSelectReceive(p.game, p.offer, p)
	}

	// On "Cancel" button
	if p.selected == numResources+1 {
		return PhaseIdle(p.game)
	}

	// Pressing confirm on a resource does nothing (use left/right)
	return p
}

func (p *phaseTradeOffer) BoardCursor() interface{} {
	return nil
}

func (p *phaseTradeOffer) Menu() string {
	var lines []string
	player := &p.game.Players[p.game.PlayerTurn]

	lines = append(lines, "What do you want to offer?")
	lines = append(lines, "")

	// Render each resource with amount
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

	lines = append(lines, "") // Empty line before buttons

	// Confirm button
	confirmLine := "Confirm"
	if p.selected == len(board.RESOURCE_TYPES) {
		confirmLine = player.Render("> ") + confirmLine
	} else {
		confirmLine = "  " + confirmLine
	}
	lines = append(lines, confirmLine)

	// Cancel button
	cancelLine := "Cancel"
	if p.selected == len(board.RESOURCE_TYPES)+1 {
		cancelLine = player.Render("> ") + cancelLine
	} else {
		cancelLine = "  " + cancelLine
	}
	lines = append(lines, cancelLine)

	return strings.Join(lines, "\n")
}

func (p *phaseTradeOffer) HelpText() string {
	if p.selected < len(board.RESOURCE_TYPES) {
		return "Use ←/→ to adjust amount, ↑/↓ to move cursor"
	}
	return "Use ↑/↓ to move cursor, Enter to confirm"
}

// phaseTradeSelectReceive allows player to select which resources to receive
// Use up/down to move cursor, left/right to adjust amounts
type phaseTradeSelectReceive struct {
	game          *Game
	offer         map[board.ResourceType]int
	request       map[board.ResourceType]int
	previousPhase Phase // The offer phase to return to on cancel
	selected      int   // 0-4 for resources, 5 for confirm, 6 for cancel
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
	numOptions := numResources + 2 // resources + confirm + cancel

	switch direction {
	case "up":
		p.selected--
		if p.selected < 0 {
			p.selected = numOptions - 1
		}
	case "down":
		p.selected++
		if p.selected >= numOptions {
			p.selected = 0
		}
	case "left":
		// Decrement resource amount if cursor on a resource
		if p.selected < numResources {
			resourceType := board.RESOURCE_TYPES[p.selected]
			if p.request[resourceType] > 0 {
				p.request[resourceType]--
			}
		}
	case "right":
		// Increment resource amount if cursor on a resource
		if p.selected < numResources {
			resourceType := board.RESOURCE_TYPES[p.selected]
			// No upper limit for request amounts
			p.request[resourceType]++
		}
	}
}

func (p *phaseTradeSelectReceive) Confirm() Phase {
	numResources := len(board.RESOURCE_TYPES)

	// On "Confirm" button
	if p.selected == numResources {
		// Check if requesting anything
		totalRequested := 0
		for _, amount := range p.request {
			totalRequested += amount
		}
		if totalRequested == 0 {
			// Stay in phase if nothing requested
			return p
		}

		// Now validate the complete trade
		return p.validateAndExecuteTrade()
	}

	// On "Cancel" button - go back to offer phase
	if p.selected == numResources+1 {
		return p.previousPhase
	}

	// Pressing confirm on a resource does nothing (use left/right)
	return p
}

func (p *phaseTradeSelectReceive) BoardCursor() interface{} {
	return nil
}

func (p *phaseTradeSelectReceive) Menu() string {
	var lines []string
	player := &p.game.Players[p.game.PlayerTurn]

	lines = append(lines, "What do you want to receive?")
	lines = append(lines, "")

	// Render each resource with amount
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

	lines = append(lines, "") // Empty line before buttons

	// Confirm button
	confirmLine := "Confirm"
	if p.selected == len(board.RESOURCE_TYPES) {
		confirmLine = player.Render("> ") + confirmLine
	} else {
		confirmLine = "  " + confirmLine
	}
	lines = append(lines, confirmLine)

	// Cancel button
	cancelLine := "Cancel"
	if p.selected == len(board.RESOURCE_TYPES)+1 {
		cancelLine = player.Render("> ") + cancelLine
	} else {
		cancelLine = "  " + cancelLine
	}
	lines = append(lines, cancelLine)

	return strings.Join(lines, "\n")
}

func (p *phaseTradeSelectReceive) HelpText() string {
	if p.selected < len(board.RESOURCE_TYPES) {
		return "Use ←/→ to adjust amount, ↑/↓ to move cursor"
	}
	return "Use ↑/↓ to move cursor, Enter to confirm"
}

// validateAndExecuteTrade checks what type of trade this is and executes it
func (p *phaseTradeSelectReceive) validateAndExecuteTrade() Phase {
	player := &p.game.Players[p.game.PlayerTurn]

	// Check if it's a valid bank trade (4 of one type → 1 of one type)
	if tradeType, offeredResource, requestedResource := p.isBankTrade(); tradeType == "bank" {
		// Validate player has enough resources
		if player.Resources[offeredResource] < 4 {
			return PhaseIdleWithNotification(p.game, "Not enough resources for bank trade!")
		}

		// Execute bank trade
		player.Resources[offeredResource] -= 4
		player.Resources[requestedResource]++

		p.game.LogAction(fmt.Sprintf("%s traded 4 %s for 1 %s with the bank",
			player.RenderName(), offeredResource, requestedResource))

		return PhaseIdleWithNotification(p.game,
			fmt.Sprintf("Traded 4 %s for 1 %s!", offeredResource, requestedResource))
	}

	// Future: Check for harbor trades here
	// if tradeType, resource, amount := p.isHarborTrade(); tradeType == "harbor" { ... }

	// Future: Check for player trades here
	// if p.isPlayerTrade() { return PhaseSelectTradePartner(...) }

	// Not a recognized trade type
	return PhaseIdleWithNotification(p.game, "Trade type not yet implemented")
}

// isBankTrade checks if the trade is exactly 4:1 (bank trade)
func (p *phaseTradeSelectReceive) isBankTrade() (string, board.ResourceType, board.ResourceType) {
	// Count offered resources
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

	// Count requested resources
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

	// Bank trade: exactly 4 of one → exactly 1 of one
	if totalOffered == 4 && offeredTypes == 1 && totalRequested == 1 && requestedTypes == 1 {
		return "bank", offeredResource, requestedResource
	}

	return "unknown", "", ""
}
