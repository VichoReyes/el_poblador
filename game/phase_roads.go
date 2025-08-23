package game

import (
	"el_poblador/board"
	"fmt"
	"strings"
)

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
	// Don't allow cancelling free roads (from development cards)
	if p.isFree {
		return p
	}
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
	// Don't allow cancelling free roads (from development cards)
	if p.isFree {
		return p
	}
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