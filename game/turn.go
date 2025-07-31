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
		panic("Play Knight not implemented")
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
		panic("Play Development Card not implemented")
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
			// TODO: Implement development card purchase
			panic("Development card purchase not implemented")
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
			return p // Invalid selection, stay in same phase
		}
	}
	return PhaseRoadEnd(p.game, p.cursorCross, p.previousPhase)
}

func (p *phaseRoadStart) Cancel() Phase {
	return p.previousPhase
}

func (p *phaseRoadStart) HelpText() string {
	return "Select the starting point for your road (must be connected to your existing road network or settlement)"
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
}

func PhaseRoadEnd(game *Game, startCross board.CrossCoord, previousPhase Phase) Phase {
	// Start with the first neighbor of the start cross
	neighbors := startCross.Neighbors()
	var cursorCross board.CrossCoord
	if len(neighbors) > 0 {
		cursorCross = neighbors[0]
	} else {
		cursorCross = startCross // fallback, though this shouldn't happen
	}

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

	// Create the path coordinate
	pathCoord := board.NewPathCoord(p.startCross, p.cursorCross)

	// Check if the road can be placed here
	if !p.game.board.CanPlaceRoad(pathCoord, playerId) {
		return p // Invalid selection, stay in same phase
	}

	// Consume resources and build the road
	if !player.BuildRoad() {
		return p // Not enough resources, stay in same phase
	}

	// Place the road on the board
	p.game.board.SetRoad(pathCoord, playerId)

	// Return to idle phase with notification
	return PhaseIdleWithNotification(p.game, "Road built!")
}

func (p *phaseRoadEnd) Cancel() Phase {
	return PhaseRoadStart(p.game, p.previousPhase)
}

func (p *phaseRoadEnd) HelpText() string {
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
		dest, ok = p.cursorCross.Up()
	case "down":
		dest, ok = p.cursorCross.Down()
	case "left":
		dest, ok = p.cursorCross.Left()
	case "right":
		dest, ok = p.cursorCross.Right()
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

	// Check if the settlement can be placed here
	if !p.game.board.CanPlaceSettlementForPlayer(p.cursorCross, playerId) {
		return p // Invalid selection, stay in same phase
	}

	// Consume resources and build the settlement
	if !player.BuildSettlement() {
		return p // Not enough resources, stay in same phase
	}

	// Place the settlement on the board
	p.game.board.SetSettlement(p.cursorCross, playerId)

	// Return to idle phase with notification
	return PhaseIdleWithNotification(p.game, "Settlement built!")
}

func (p *phaseSettlementPlacement) Cancel() Phase {
	return p.previousPhase
}

func (p *phaseSettlementPlacement) HelpText() string {
	return "Select where to place your settlement (must be connected to your road network)"
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

	// Check if the city can be placed here (upgrade existing settlement)
	if !p.game.board.CanUpgradeToCity(p.cursorCross, playerId) {
		return p // Invalid selection, stay in same phase
	}

	// Consume resources and build the city
	if !player.BuildCity() {
		return p // Not enough resources, stay in same phase
	}

	// Upgrade the settlement to a city on the board
	p.game.board.UpgradeToCity(p.cursorCross, playerId)

	// Return to idle phase with notification
	return PhaseIdleWithNotification(p.game, "City built!")
}

func (p *phaseCityPlacement) Cancel() Phase {
	return p.previousPhase
}

func (p *phaseCityPlacement) HelpText() string {
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
