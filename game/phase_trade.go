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

	return "unknown", 0, 0
}

// phaseTradeMenu is a phase that allows the player to select the type of trade they want to do
type phaseTradeMenu struct {
	phaseWithOptions
	previous Phase // Moved previous here
}

func PhaseTradeMenu(game *Game, previous Phase) Phase {
	return &phaseTradeMenu{
		phaseWithOptions: phaseWithOptions{
			game:    game,
			options: []string{"Trade with Bank/Harbors", "Trade with Players"},
		},
		previous: previous, // Assigned previous here
	}
}

func (p *phaseTradeMenu) Confirm() Phase {
	switch p.selected {
	case 0:
		return PhaseTradeOffer(p.game)
	case 1:
		return PhasePlayerTrade(p.game, p)
	}
	return p
}

func (p *phaseTradeMenu) Cancel() Phase { // Added Cancel method
	return p.previous
}

func (p *phaseTradeMenu) HelpText() string {
	return "Select the type of trade you want to perform"
}

// phasePlayerTrade is the main menu for player-to-player trading.
// It displays incoming and outgoing offers and allows creating new ones.
type phasePlayerTrade struct {
	game     *Game
	previous Phase
	selected int
	// combined list of offers for easier selection
	offers []*TradeOffer
}

func PhasePlayerTrade(game *Game, previous Phase) Phase {
	p := &phasePlayerTrade{
		game:     game,
		previous: previous,
		selected: 0,
	}
	p.updateOffers()
	return p
}

func (p *phasePlayerTrade) updateOffers() {
	p.offers = []*TradeOffer{}
	playerID := p.game.PlayerTurn

	// Incoming
	for i := range p.game.TradeOffers {
		offer := &p.game.TradeOffers[i]
		if offer.Status == OfferIsPending && (offer.TargetID == playerID || offer.TargetID == -1) && offer.OffererID != playerID {
			p.offers = append(p.offers, offer)
		}
	}
	// Outgoing
	for i := range p.game.TradeOffers {
		offer := &p.game.TradeOffers[i]
		if offer.Status == OfferIsPending && offer.OffererID == playerID {
			p.offers = append(p.offers, offer)
		}
	}
}

func (p *phasePlayerTrade) MoveCursor(direction string) {
	// +1 for the "Create New Offer" option
	totalOptions := len(p.offers) + 1

	switch direction {
	case "up":
		p.selected--
		if p.selected < 0 {
			p.selected = totalOptions - 1
		}
	case "down":
		p.selected++
		if p.selected >= totalOptions {
			p.selected = 0
		}
	}
}

func (p *phasePlayerTrade) Confirm() Phase {
	totalOffers := len(p.offers)

	if p.selected < totalOffers {
		selectedOffer := p.offers[p.selected]
		playerID := p.game.PlayerTurn

		if selectedOffer.OffererID == playerID {
			// Retract the offer
			selectedOffer.Status = OfferIsRetracted
			p.game.LogAction(fmt.Sprintf("%s retracted an offer.", p.game.Players[playerID].RenderName()))
		} else {
			// Accept the offer
			return p.acceptOffer(selectedOffer)
		}
		p.updateOffers() // Refresh the list
		return p
	}

	// "Create New Offer" was selected
	return PhaseCreateOffer(p.game, p)
}

func (p *phasePlayerTrade) acceptOffer(offer *TradeOffer) Phase {
	accepter := &p.game.Players[p.game.PlayerTurn]
	offerer := &p.game.Players[offer.OffererID]

	// 1. Validate the trade (concrete trades only)
	if !p.canExecuteTrade(offerer, accepter, offer) {
		// Cannot accept ambiguous trades directly, must be countered
		// Also handles insufficient resources
		return PhaseIdleWithNotification(p.game, "Cannot execute this trade (not enough resources or offer is ambiguous).")
	}

	// 2. Exchange resources
	for resource, amount := range offer.Offering {
		offerer.Resources[resource] -= amount
		accepter.Resources[resource] += amount
	}
	for resource, amount := range offer.Requesting {
		accepter.Resources[resource] -= amount
		offerer.Resources[resource] += amount
	}

	// 3. Update offer status
	offer.Status = OfferIsCompleted

	// 4. Log the action
	p.game.LogAction(fmt.Sprintf("%s accepted a trade from %s.", accepter.RenderName(), offerer.RenderName()))

	// 5. Return to the idle phase with a notification
	return PhaseIdleWithNotification(p.game, "Trade successful!")
}

// canExecuteTrade checks if a trade is concrete and if both players have the resources
func (p *phasePlayerTrade) canExecuteTrade(offerer, accepter *Player, offer *TradeOffer) bool {
	// Check for ambiguous resources
	for res := range offer.Offering {
		if res == board.ResourceInvalid {
			return false
		}
	}
	for res := range offer.Requesting {
		if res == board.ResourceInvalid {
			return false
		}
	}

	// Check if offerer has the resources
	for resource, amount := range offer.Offering {
		if offerer.Resources[resource] < amount {
			return false
		}
	}

	// Check if accepter has the resources
	for resource, amount := range offer.Requesting {
		if accepter.Resources[resource] < amount {
			return false
		}
	}

	return true
}

func (p *phasePlayerTrade) Cancel() Phase {
	return p.previous
}

func (p *phasePlayerTrade) BoardCursor() interface{} {
	return nil
}

func (p *phasePlayerTrade) Menu() string {
	var lines []string
	player := &p.game.Players[p.game.PlayerTurn]
	playerID := p.game.PlayerTurn

	lines = append(lines, "Player Trading")
	lines = append(lines, "")

	currentIndex := 0

	if len(p.offers) == 0 {
		lines = append(lines, "  (No pending offers)")
	}

	for _, offer := range p.offers {
		var tag string
		if offer.OffererID == playerID {
			tag = "[My Offer]"
		} else {
			tag = "[Incoming]"
		}

		line := fmt.Sprintf("  %s %s", tag, formatOffer(*offer, p.game))
		if currentIndex == p.selected {
			line = player.Render("> ") + line
		}
		lines = append(lines, line)
		currentIndex++
	}
	lines = append(lines, "")

	// Display Create New Offer
	createLine := "Create New Offer"
	if currentIndex == p.selected {
		createLine = player.Render("> ") + createLine
	}
	lines = append(lines, createLine)

	return strings.Join(lines, "\n")
}

func (p *phasePlayerTrade) HelpText() string {
	return "Use ↑/↓ to navigate, Enter to select, Esc to go back."
}

// formatOffer is a helper to create a readable string for a TradeOffer
func formatOffer(offer TradeOffer, g *Game) string {
	offererName := g.Players[offer.OffererID].RenderName()
	offeringStr := formatResourceMap(offer.Offering)
	requestingStr := formatResourceMap(offer.Requesting)

	var targetStr string
	if offer.TargetID == -1 {
		targetStr = "all players"
	} else {
		targetStr = g.Players[offer.TargetID].RenderName()
	}

	return fmt.Sprintf("%s offers %s for %s to %s", offererName, offeringStr, requestingStr, targetStr)
}

// formatResourceMap is a helper to create a readable string from a resource map
func formatResourceMap(resources map[board.ResourceType]int) string {
	if len(resources) == 0 {
		return "anything"
	}
	var parts []string
	for res, amount := range resources {
		if res == board.ResourceInvalid {
			parts = append(parts, fmt.Sprintf("%d other", amount))
		} else {
			parts = append(parts, fmt.Sprintf("%d %s", amount, res))
		}
	}
	return strings.Join(parts, ", ")
}

// phaseCreateOfferOffer is the first step in creating a player-to-player trade offer.
// It allows the player to select the resources they want to offer.
type phaseCreateOfferOffer struct {
	game             *Game
	playerTradePhase Phase // Store the original phasePlayerTrade here
	previous         Phase // For Cancel() to return to the immediate previous phase
	offer            map[board.ResourceType]int
	selected         int
}

func PhaseCreateOffer(game *Game, playerTradePhase Phase) Phase { // playerTradePhase is the previous
	offer := make(map[board.ResourceType]int)
	// Include ResourceInvalid for ambiguous offers
	for _, resourceType := range append(board.RESOURCE_TYPES, board.ResourceInvalid) {
		offer[resourceType] = 0
	}

	return &phaseCreateOfferOffer{
		game:             game,
		playerTradePhase: playerTradePhase, // Assign here
		previous:         playerTradePhase, // Also set previous for direct cancel
		offer:            offer,
		selected:         0,
	}
}

func (p *phaseCreateOfferOffer) MoveCursor(direction string) {
	numOptions := len(board.RESOURCE_TYPES) + 1 // +1 for ResourceInvalid
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
		resourceType := p.getSelectedResourceType()
		if p.offer[resourceType] > 0 {
			p.offer[resourceType]--
		}
	case "right":
		resourceType := p.getSelectedResourceType()
		player := &p.game.Players[p.game.PlayerTurn]
		// For non-invalid resources, check against player's inventory
		if resourceType != board.ResourceInvalid {
			maxAvailable := player.Resources[resourceType]
			if p.offer[resourceType] < maxAvailable {
				p.offer[resourceType]++
			}
		} else { // For ResourceInvalid, allow any amount
			p.offer[resourceType]++
		}
	}
}

func (p *phaseCreateOfferOffer) getSelectedResourceType() board.ResourceType {
	if p.selected < len(board.RESOURCE_TYPES) {
		return board.RESOURCE_TYPES[p.selected]
	}
	return board.ResourceInvalid
}

func (p *phaseCreateOfferOffer) Confirm() Phase {
	totalOffered := 0
	for _, amount := range p.offer {
		totalOffered += amount
	}
	if totalOffered == 0 {
		return p // Don't proceed with an empty offer
	}
	return PhaseCreateOfferRequest(p.game, p.playerTradePhase, p, p.offer)
}

func (p *phaseCreateOfferOffer) Cancel() Phase {
	return p.previous
}

func (p *phaseCreateOfferOffer) BoardCursor() interface{} {
	return nil
}

func (p *phaseCreateOfferOffer) Menu() string {
	var lines []string
	player := &p.game.Players[p.game.PlayerTurn]
	resourceOptions := append(board.RESOURCE_TYPES, board.ResourceInvalid)

	lines = append(lines, "What will you offer?")
	lines = append(lines, "")

	for i, resourceType := range resourceOptions {
		amount := p.offer[resourceType]
		var line string
		if resourceType == board.ResourceInvalid {
			line = fmt.Sprintf("Other: %d", amount)
		} else {
			maxAvailable := player.Resources[resourceType]
			line = fmt.Sprintf("%s: %d / %d", resourceType, amount, maxAvailable)
		}

		if i == p.selected {
			line = player.Render("> ") + line
		} else {
			line = "  " + line
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func (p *phaseCreateOfferOffer) HelpText() string {
	return "Use ←/→ to adjust, ↑/↓ to move, Enter to confirm, Esc to cancel"
}

// phaseCreateOfferRequest is the second step in creating a player-to-player trade offer.
// It allows the player to select the resources they want to request.
type phaseCreateOfferRequest struct {
	game             *Game
	playerTradePhase Phase // Store the original phasePlayerTrade here
	previous         Phase // Keep previous for Cancel() to return to phaseCreateOfferOffer
	offer            map[board.ResourceType]int
	request          map[board.ResourceType]int
	selected         int
}

func PhaseCreateOfferRequest(game *Game, playerTradePhase, previous Phase, offer map[board.ResourceType]int) Phase {
	request := make(map[board.ResourceType]int)
	// Include ResourceInvalid for ambiguous offers
	for _, resourceType := range append(board.RESOURCE_TYPES, board.ResourceInvalid) {
		request[resourceType] = 0
	}

	return &phaseCreateOfferRequest{
		game:             game,
		playerTradePhase: playerTradePhase, // Assign here
		previous:         previous,
		offer:            offer,
		request:          request,
		selected:         0,
	}
}

func (p *phaseCreateOfferRequest) MoveCursor(direction string) {
	numOptions := len(board.RESOURCE_TYPES) + 1 // +1 for ResourceInvalid
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
		resourceType := p.getSelectedResourceType()
		if p.request[resourceType] > 0 {
			p.request[resourceType]--
		}
	case "right":
		resourceType := p.getSelectedResourceType()
		p.request[resourceType]++
	}
}

func (p *phaseCreateOfferRequest) getSelectedResourceType() board.ResourceType {
	if p.selected < len(board.RESOURCE_TYPES) {
		return board.RESOURCE_TYPES[p.selected]
	}
	return board.ResourceInvalid
}

func (p *phaseCreateOfferRequest) Confirm() Phase {
	totalRequested := 0
	for _, amount := range p.request {
		totalRequested += amount
	}
	if totalRequested == 0 {
		return p // Don't proceed with an empty request
	}
	return PhaseCreateOfferTarget(p.game, p.playerTradePhase, p, p.offer, p.request)
}

func (p *phaseCreateOfferRequest) Cancel() Phase {
	return p.previous
}

func (p *phaseCreateOfferRequest) BoardCursor() interface{} {
	return nil
}

func (p *phaseCreateOfferRequest) Menu() string {
	var lines []string
	player := &p.game.Players[p.game.PlayerTurn]
	resourceOptions := append(board.RESOURCE_TYPES, board.ResourceInvalid)

	lines = append(lines, "What will you request?")
	lines = append(lines, "")

	for i, resourceType := range resourceOptions {
		amount := p.request[resourceType]
		var line string
		if resourceType == board.ResourceInvalid {
			line = fmt.Sprintf("Other: %d", amount)
		} else {
			line = fmt.Sprintf("%s: %d", resourceType, amount)
		}

		if i == p.selected {
			line = player.Render("> ") + line
		} else {
			line = "  " + line
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func (p *phaseCreateOfferRequest) HelpText() string {
	return "Use ←/→ to adjust, ↑/↓ to move, Enter to confirm, Esc to cancel"
}

// phaseCreateOfferTarget is the final step in creating a player-to-player trade offer.
// It allows the player to select the target for the offer.
type phaseCreateOfferTarget struct {
	game             *Game
	playerTradePhase Phase // Store the original phasePlayerTrade here
	previous         Phase // Keep previous for Cancel() to return to phaseCreateOfferRequest
	offer            map[board.ResourceType]int
	request          map[board.ResourceType]int
	selected         int
}

func PhaseCreateOfferTarget(game *Game, playerTradePhase, previous Phase, offer, request map[board.ResourceType]int) Phase {
	return &phaseCreateOfferTarget{
		game:             game,
		playerTradePhase: playerTradePhase, // Assign here
		previous:         previous,
		offer:            offer,
		request:          request,
		selected:         0,
	}
}

func (p *phaseCreateOfferTarget) getTargetOptions() []string {
	options := []string{"All Players"}
	for i, player := range p.game.Players {
		if i != p.game.PlayerTurn {
			options = append(options, player.Name)
		}
	}
	return options
}

func (p *phaseCreateOfferTarget) MoveCursor(direction string) {
	numOptions := len(p.getTargetOptions())
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
	}
}

func (p *phaseCreateOfferTarget) Confirm() Phase {
	options := p.getTargetOptions()
	selectedOption := options[p.selected]
	targetID := -1 // Default to "All Players"

	if selectedOption != "All Players" {
		for i, player := range p.game.Players {
			if player.Name == selectedOption {
				targetID = i
				break
			}
		}
	}

	// Create and add the offer
	newOffer := TradeOffer{
		ID:          len(p.game.TradeOffers), // Simple unique ID for now
		OffererID:   p.game.PlayerTurn,
		TargetID:    targetID,
		Offering:    p.offer,
		Requesting:  p.request,
		Status:      OfferIsPending,
		InReplyToID: 0, // Not a reply
	}
	p.game.TradeOffers = append(p.game.TradeOffers, newOffer)
	p.game.LogAction(fmt.Sprintf("%s made a trade offer.", p.game.Players[p.game.PlayerTurn].RenderName()))

	// Return to the main player trade menu
	return p.playerTradePhase
}

func (p *phaseCreateOfferTarget) Cancel() Phase {
	return p.previous
}

func (p *phaseCreateOfferTarget) BoardCursor() interface{} {
	return nil
}

func (p *phaseCreateOfferTarget) Menu() string {
	var lines []string
	player := &p.game.Players[p.game.PlayerTurn]
	options := p.getTargetOptions()

	lines = append(lines, "Who do you want to trade with?")
	lines = append(lines, "")

	for i, option := range options {
		line := option
		if i == p.selected {
			line = player.Render("> ") + line
		} else {
			line = "  " + line
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func (p *phaseCreateOfferTarget) HelpText() string {
	return "Use ↑/↓ to select, Enter to confirm, Esc to go back."
}

