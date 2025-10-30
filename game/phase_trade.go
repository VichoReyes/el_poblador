package game

import (
	"el_poblador/board"
	"fmt"
	"strings"
)

// phaseTradeOffer allows player to specify resources to offer
// Use up/down to move cursor, left/right to adjust amounts
type phaseTradeOffer struct {
	game          *Game
	previousPhase Phase
	offer         map[board.ResourceType]int
	selected      int // 0-4 for resources, 5 for confirm, 6 for cancel
}

func PhaseTradeOffer(game *Game, previousPhase Phase) Phase {
	offer := make(map[board.ResourceType]int)
	for _, resourceType := range board.RESOURCE_TYPES {
		offer[resourceType] = 0
	}

	return &phaseTradeOffer{
		game:          game,
		previousPhase: previousPhase,
		offer:         offer,
		selected:      0,
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

	// On "Confirm Trade" button
	if p.selected == numResources {
		if p.isValidBankTrade() {
			return PhaseTradeSelectReceive(p.game, p.offer, p.previousPhase)
		}
		// Invalid trade - stay in phase (error shown in HelpText)
		return p
	}

	// On "Cancel" button
	if p.selected == numResources+1 {
		return p.previousPhase
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
	confirmLine := "Confirm Trade"
	if !p.isValidBankTrade() {
		confirmLine += " (need 4 of one resource)"
	}
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

// isValidBankTrade checks if offer is valid for bank trading
// Bank trade requires exactly 4 of exactly one resource type
func (p *phaseTradeOffer) isValidBankTrade() bool {
	player := &p.game.Players[p.game.PlayerTurn]
	totalOffered := 0
	resourcesOffered := 0
	var offeredResource board.ResourceType

	for resourceType, amount := range p.offer {
		if amount > 0 {
			resourcesOffered++
			offeredResource = resourceType
			totalOffered += amount
		}
	}

	// Must offer exactly 4 of exactly one resource
	if totalOffered != 4 || resourcesOffered != 1 {
		return false
	}

	// Must have enough resources
	return player.Resources[offeredResource] >= 4
}

// phaseTradeSelectReceive allows player to select which resource to receive
type phaseTradeSelectReceive struct {
	phaseWithOptions
	offer         map[board.ResourceType]int
	previousPhase Phase // The original idle phase to return to
}

func PhaseTradeSelectReceive(game *Game, offer map[board.ResourceType]int, previousPhase Phase) Phase {
	// Build list of resources that can be received (all except what's being offered)
	var options []string
	for _, resourceType := range board.RESOURCE_TYPES {
		if offer[resourceType] == 0 {
			options = append(options, string(resourceType))
		}
	}
	options = append(options, "Cancel")

	return &phaseTradeSelectReceive{
		phaseWithOptions: phaseWithOptions{
			game:    game,
			options: options,
		},
		offer:         offer,
		previousPhase: previousPhase,
	}
}

func (p *phaseTradeSelectReceive) Confirm() Phase {
	// Check if cancel was selected
	if p.selected == len(p.options)-1 {
		// Go back to offer phase with current offer preserved
		phase := PhaseTradeOffer(p.game, p.previousPhase).(*phaseTradeOffer)
		phase.offer = p.offer
		return phase
	}

	// Execute the trade
	selectedResource := board.ResourceType(p.options[p.selected])
	player := &p.game.Players[p.game.PlayerTurn]

	// Find what resource is being offered
	var offeredResource board.ResourceType
	var offeredAmount int
	for resourceType, amount := range p.offer {
		if amount > 0 {
			offeredResource = resourceType
			offeredAmount = amount
			break
		}
	}

	// Deduct offered resources
	player.Resources[offeredResource] -= offeredAmount

	// Add received resource
	player.Resources[selectedResource]++

	// Log the trade
	p.game.LogAction(fmt.Sprintf("%s traded %d %s for 1 %s with the bank",
		player.RenderName(), offeredAmount, offeredResource, selectedResource))

	return PhaseIdleWithNotification(p.game,
		fmt.Sprintf("Traded %d %s for 1 %s!", offeredAmount, offeredResource, selectedResource))
}

func (p *phaseTradeSelectReceive) HelpText() string {
	return "Select a resource to receive"
}
