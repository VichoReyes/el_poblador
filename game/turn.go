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
		if player.HasKnightCard() {
			// TODO: Implement robber mechanics - for now just play the card
			if player.PlayDevCard(DevCardKnight) {
				// Return to idle phase with notification
				return PhaseIdleWithNotification(p.game, "Knight card played! Robber mechanics coming soon...")
			}
		}
		// If no knight card available, stay in same phase
		return p
	default:
		panic("Invalid option selected")
	}
}

func rollDice(game *Game) Phase {
	game.lastDice = [2]int{rand.IntN(6) + 1, rand.IntN(6) + 1}
	sum := game.lastDice[0] + game.lastDice[1]
	if sum == 7 {
		// TODO: discarding of > 7 cards
		// also TODO: implement robber
		// for now go to idle phase
		return PhaseIdle(game)
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
}

func PhaseRoadStart(game *Game, previousPhase Phase) Phase {
	cursorCross := game.board.ValidCrossCoord()
	return &phaseRoadStart{
		game:          game,
		cursorCross:   cursorCross,
		previousPhase: previousPhase,
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
	return PhaseRoadEnd(p.game, p.cursorCross, p.previousPhase)
}

func (p *phaseRoadStart) Cancel() Phase {
	return p.previousPhase
}

func (p *phaseRoadStart) HelpText() string {
	if p.invalid != "" {
		return p.invalid
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
}

func PhaseRoadEnd(game *Game, startCross board.CrossCoord, previousPhase Phase) Phase {
	// Start with the first neighbor of the start cross
	neighbors := startCross.Neighbors()
	if len(neighbors) == 0 {
		panic("No neighbors found for start cross")
	}
	cursorCross := neighbors[0]

	return &phaseRoadEnd{
		game:          game,
		startCross:    startCross,
		cursorCross:   cursorCross,
		previousPhase: previousPhase,
	}
}

func (p *phaseRoadEnd) Confirm() Phase {
	player := p.game.players[p.game.playerTurn]
	playerId := p.game.playerTurn
	pathCoord := board.NewPathCoord(p.startCross, p.cursorCross)

	if !p.game.board.CanPlaceRoad(pathCoord, playerId) {
		p.invalid = "Can't build road here"
		return p
	}

	if !player.BuildRoad() {
		p.invalid = "Not enough resources"
		return p
	}

	p.game.board.SetRoad(pathCoord, playerId)
	return PhaseIdleWithNotification(p.game, "Road built!")
}

func (p *phaseRoadEnd) Cancel() Phase {
	return PhaseRoadStart(p.game, p.previousPhase)
}

func (p *phaseRoadEnd) HelpText() string {
	if p.invalid != "" {
		return p.invalid
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

	// Add available development card options
	if player.HasKnightCard() {
		options = append(options, "Knight")
	}
	// TODO: Add other development card types as they're implemented

	// Always add cancel option
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
	player := p.game.players[p.game.playerTurn]

	switch p.selected {
	case 0: // Knight (if available)
		if player.HasKnightCard() {
			if player.PlayDevCard(DevCardKnight) {
				return PhaseIdleWithNotification(p.game, "Knight card played! Robber mechanics coming soon...")
			}
		}
		return p
	default: // Cancel or invalid
		return p.previousPhase
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
