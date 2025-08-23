package game

import (
	"el_poblador/board"
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
		// Play Knight card
		player := p.game.players[p.game.playerTurn]
		if player.PlayDevCard(DevCardKnight) {
			return PhasePlaceRobber(p.game, p)
		}
		return p
	default:
		panic("Invalid option selected")
	}
}

type phasePlaceRobber struct {
	game         *Game
	tileCoord    board.TileCoord
	continuation Phase
	invalid      string
}

func PhasePlaceRobber(game *Game, continuation Phase) Phase {
	return &phasePlaceRobber{
		game:         game,
		continuation: continuation,
	}
}

func (p *phasePlaceRobber) BoardCursor() interface{} {
	return p.tileCoord
}

func (p *phasePlaceRobber) MoveCursor(direction string) {
	dest, ok := moveTileCursor(p.tileCoord, direction)
	if !ok {
		return
	}
	p.tileCoord = dest
}

func moveTileCursor(tileCoord board.TileCoord, direction string) (board.TileCoord, bool) {
	switch direction {
	case "up":
		return tileCoord.Up()
	case "down":
		return tileCoord.Down()
	case "left":
		return tileCoord.Left()
	case "right":
		return tileCoord.Right()
	default:
		return tileCoord, false
	}
}

func (p *phasePlaceRobber) HelpText() string {
	if p.invalid != "" {
		return p.invalid
	}
	return "Select a tile to place the robber on"
}

func (p *phasePlaceRobber) Confirm() Phase {
	// Check if trying to place robber on the same tile it's already on
	if p.tileCoord == p.game.board.GetRobber() {
		p.invalid = "Robber cannot be moved to the same tile it's already on"
		return p
	}
	
	playerIds := p.game.board.PlaceRobber(p.tileCoord)
	var stealablePlayers []Player
	for _, playerId := range playerIds {
		p := p.game.players[playerId]
		if p.TotalResources() > 0 {
			stealablePlayers = append(stealablePlayers, p)
		}
	}
	if len(stealablePlayers) == 0 { // no one to steal from? skip
		return p.continuation
	}
	return &phaseStealCard{
		game:             p.game,
		continuation:     p.continuation,
		stealablePlayers: stealablePlayers,
	}
}

type phaseStealCard struct {
	game             *Game
	continuation     Phase
	stealablePlayers []Player
	selected         int
}

func (p *phaseStealCard) BoardCursor() interface{} {
	return nil
}

func (p *phaseStealCard) MoveCursor(direction string) {
	switch direction {
	case "up":
		p.selected--
	case "down":
		p.selected++
	}
	p.selected = (p.selected + len(p.stealablePlayers)) % len(p.stealablePlayers)
}

func (p *phaseStealCard) HelpText() string {
	return "Select a player to steal from"
}

func (p *phaseStealCard) Confirm() Phase {
	player := p.stealablePlayers[p.selected]
	var resourcePool []board.ResourceType
	for resType, count := range player.resources {
		for i := 0; i < count; i++ {
			resourcePool = append(resourcePool, resType)
		}
	}
	if len(resourcePool) > 0 {
		selectedResource := resourcePool[rand.IntN(len(resourcePool))]
		player.resources[selectedResource] -= 1
		p.game.players[p.game.playerTurn].AddResource(selectedResource)
	}
	return p.continuation
}

func (p *phaseStealCard) Menu() string {
	var paddedOptions []string
	for i, player := range p.stealablePlayers {
		if i == p.selected {
			paddedOptions = append(paddedOptions, "> "+player.Render(player.Name))
		} else {
			paddedOptions = append(paddedOptions, player.Render(" "+player.Name))
		}
	}
	return strings.Join(paddedOptions, "\n")
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

func (p *phaseDiceRoll) HelpText() string {
	return "Time to roll the dice"
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
		panic("Trade not implemented")
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
			return PhaseRoadStart(p.game, p)
		}
		return p
	case 1: // Settlement
		if player.CanBuildSettlement() {
			return PhaseSettlementPlacement(p.game, p)
		}
		return p
	case 2: // City
		if player.CanBuildCity() {
			return PhaseCityPlacement(p.game, p)
		}
		return p
	case 3: // Development Card
		if player.CanBuyDevelopmentCard() {
			if player.BuyDevelopmentCard() {
				if card := p.game.DrawDevelopmentCard(); card != nil {
					player.hiddenDevCards = append(player.hiddenDevCards, *card)
					return p.previousPhase
				}
			}
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

type phaseRoadStart struct {
	game          *Game
	cursorCross   board.CrossCoord
	previousPhase Phase
	invalid       string
	isFree        bool
	continuation  Phase
	helpPrefix    string
}

// Phase for building a road by paying for it
func PhaseRoadStart(game *Game, previousPhase Phase) Phase {
	return newPhaseRoadStart(game, previousPhase, false, nil, "")
}

// Phase for building two roads by using a development card
func PhaseRoadBuilding(game *Game) Phase {
	// First free road - continuation will be second free road
	secondRoadPhase := newPhaseRoadStart(game, PhaseIdle(game), true, PhaseIdleWithNotification(game, "Two free roads built!"), "second free")
	return newPhaseRoadStart(game, PhaseIdle(game), true, secondRoadPhase, "first free")
}

func newPhaseRoadStart(game *Game, previousPhase Phase, isFree bool, continuation Phase, helpPrefix string) Phase {
	cursorCross := game.board.ValidCrossCoord()
	finalContinuation := continuation
	if finalContinuation == nil {
		finalContinuation = previousPhase
	}
	return &phaseRoadStart{
		game:          game,
		cursorCross:   cursorCross,
		previousPhase: previousPhase,
		isFree:        isFree,
		continuation:  finalContinuation,
		helpPrefix:    helpPrefix,
	}
}

func (p *phaseRoadStart) Confirm() Phase {
	// Check if player has a road or settlement connected to this crossing
	playerId := p.game.playerTurn
	if !p.game.board.HasRoadConnected(p.cursorCross, playerId) {
		// Check if player has a settlement at this crossing
		if !p.game.board.HasSettlementAt(p.cursorCross, playerId) {
			p.invalid = "You must have a road or settlement connected to this crossing"
			return p // Invalid selection, stay in same phase
		}
	}
	return newPhaseRoadEnd(p.game, p.cursorCross, p.previousPhase, p.isFree, p.continuation, p.helpPrefix)
}

func (p *phaseRoadStart) Cancel() Phase {
	return p.previousPhase
}

func (p *phaseRoadStart) HelpText() string {
	if p.invalid != "" {
		return p.invalid
	}
	if p.helpPrefix != "" {
		return fmt.Sprintf("Select the starting point for your %s road", p.helpPrefix)
	}
	return "Select the starting point for your road"
}

func (p *phaseRoadStart) BoardCursor() interface{} {
	return p.cursorCross
}

func (p *phaseRoadStart) MoveCursor(direction string) {
	dest, ok := moveCrossCursor(p.cursorCross, direction)
	if !ok {
		return
	}
	p.cursorCross = dest
}

type phaseRoadEnd struct {
	game          *Game
	startCross    board.CrossCoord
	cursorCross   board.CrossCoord
	previousPhase Phase
	invalid       string
	isFree        bool
	continuation  Phase
	helpPrefix    string
}

func PhaseRoadEnd(game *Game, startCross board.CrossCoord, previousPhase Phase) Phase {
	return newPhaseRoadEnd(game, startCross, previousPhase, false, nil, "")
}

func newPhaseRoadEnd(game *Game, startCross board.CrossCoord, previousPhase Phase, isFree bool, continuation Phase, helpPrefix string) Phase {
	// Start with the first neighbor of the start cross
	neighbors := startCross.Neighbors()
	if len(neighbors) == 0 {
		panic("No neighbors found for start cross")
	}
	cursorCross := neighbors[0]
	
	finalContinuation := continuation
	if finalContinuation == nil {
		finalContinuation = previousPhase
	}

	return &phaseRoadEnd{
		game:          game,
		startCross:    startCross,
		cursorCross:   cursorCross,
		previousPhase: previousPhase,
		isFree:        isFree,
		continuation:  finalContinuation,
		helpPrefix:    helpPrefix,
	}
}

func (p *phaseRoadEnd) Confirm() Phase {
	player := &p.game.players[p.game.playerTurn]
	playerId := p.game.playerTurn
	pathCoord := board.NewPathCoord(p.startCross, p.cursorCross)

	if !p.game.board.CanPlaceRoad(pathCoord, playerId) {
		p.invalid = "Can't build road here"
		return p
	}

	if !p.isFree {
		if !player.BuildRoad() {
			p.invalid = "Not enough resources"
			return p
		}
	}

	p.game.board.SetRoad(pathCoord, playerId)
	
	message := "Road built!"
	if p.isFree {
		if p.helpPrefix != "" {
			message = fmt.Sprintf("%s road built!", strings.Title(p.helpPrefix))
		} else {
			message = "Free road built!"
		}
	}
	
	if p.continuation == nil {
		return PhaseIdleWithNotification(p.game, message)
	} else {
		return p.continuation
	}
}

func (p *phaseRoadEnd) Cancel() Phase {
	return newPhaseRoadStart(p.game, p.previousPhase, p.isFree, p.continuation, p.helpPrefix)
}

func (p *phaseRoadEnd) HelpText() string {
	if p.invalid != "" {
		return p.invalid
	}
	if p.helpPrefix != "" {
		return fmt.Sprintf("Select the ending point for your %s road", p.helpPrefix)
	}
	return "Select the ending point for your road"
}

func (p *phaseRoadEnd) BoardCursor() interface{} {
	return p.cursorCross
}

func (p *phaseRoadEnd) MoveCursor(direction string) {
	var dest board.CrossCoord
	var ok bool

	switch direction {
	case "up":
		dest, ok = p.startCross.Up()
	case "down":
		dest, ok = p.startCross.Down()
	case "left":
		dest, ok = p.startCross.Left()
	case "right":
		dest, ok = p.startCross.Right()
	default:
		return
	}

	if ok {
		p.cursorCross = dest
	}
}

type phaseSettlementPlacement struct {
	game          *Game
	cursorCross   board.CrossCoord
	previousPhase Phase
	invalid       string
}

func PhaseSettlementPlacement(game *Game, previousPhase Phase) Phase {
	cursorCross := game.board.ValidCrossCoord()
	return &phaseSettlementPlacement{
		game:          game,
		cursorCross:   cursorCross,
		previousPhase: previousPhase,
	}
}

func (p *phaseSettlementPlacement) Confirm() Phase {
	player := p.game.players[p.game.playerTurn]
	playerId := p.game.playerTurn

	if !p.game.board.CanPlaceSettlementForPlayer(p.cursorCross, playerId) {
		p.invalid = "Can't build settlement here"
		return p
	}

	if !player.BuildSettlement() {
		p.invalid = "Not enough resources"
		return p
	}

	p.game.board.SetSettlement(p.cursorCross, playerId)

	return PhaseIdleWithNotification(p.game, "Settlement built!")
}

func (p *phaseSettlementPlacement) Cancel() Phase {
	return p.previousPhase
}

func (p *phaseSettlementPlacement) HelpText() string {
	if p.invalid != "" {
		return p.invalid
	}
	return "Select where to place your settlement"
}

func (p *phaseSettlementPlacement) BoardCursor() interface{} {
	return p.cursorCross
}

func (p *phaseSettlementPlacement) MoveCursor(direction string) {
	dest, ok := moveCrossCursor(p.cursorCross, direction)
	if !ok {
		return
	}
	p.cursorCross = dest
}

type phaseCityPlacement struct {
	game          *Game
	cursorCross   board.CrossCoord
	previousPhase Phase
	invalid       string
}

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

func PhaseCityPlacement(game *Game, previousPhase Phase) Phase {
	cursorCross := game.board.ValidCrossCoord()
	return &phaseCityPlacement{
		game:          game,
		cursorCross:   cursorCross,
		previousPhase: previousPhase,
	}
}

func (p *phaseCityPlacement) Confirm() Phase {
	player := p.game.players[p.game.playerTurn]
	playerId := p.game.playerTurn

	if !p.game.board.CanUpgradeToCity(p.cursorCross, playerId) {
		p.invalid = "Can't upgrade to city here"
		return p
	}

	if !player.BuildCity() {
		p.invalid = "Not enough resources"
		return p
	}

	p.game.board.UpgradeToCity(p.cursorCross, playerId)

	return PhaseIdleWithNotification(p.game, "City built!")
}

func (p *phaseCityPlacement) Cancel() Phase {
	return p.previousPhase
}

func (p *phaseCityPlacement) HelpText() string {
	if p.invalid != "" {
		return p.invalid
	}
	return "Select a settlement to upgrade to a city"
}

func (p *phaseCityPlacement) BoardCursor() interface{} {
	return p.cursorCross
}

func (p *phaseCityPlacement) MoveCursor(direction string) {
	dest, ok := moveCrossCursor(p.cursorCross, direction)
	if !ok {
		return
	}
	p.cursorCross = dest
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
	previousPhase   Phase
	selectedCount   int
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

