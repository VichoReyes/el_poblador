package game

import (
	"el_poblador/board"
	"fmt"
)

type phaseInitialSettlements struct {
	game        *Game
	playerTurn  int
	cursorCross board.CrossCoord
}

func PhaseInitialSettlements(game *Game, playerTurn int) Phase {
	cursorCross := game.board.ValidCrossCoord()
	return &phaseInitialSettlements{game: game, playerTurn: playerTurn, cursorCross: cursorCross}
}

func (p *phaseInitialSettlements) Confirm() Phase {
	p.game.board.SetSettlement(p.cursorCross, p.game.players[p.playerTurn].Name)
	return PhaseInitialRoad(p.game, p.playerTurn, p.cursorCross)
}

func (p *phaseInitialSettlements) Cancel() Phase {
	return p
}

func (p *phaseInitialSettlements) HelpText() string {
	return fmt.Sprintf("%s's turn. Place your initial settlement on the board with 'enter'.", p.game.players[p.playerTurn].Name)
}

func (p *phaseInitialSettlements) CurrentCursor() interface{} {
	return p.cursorCross
}

func (p *phaseInitialSettlements) MoveCursor(direction string) {
	dest, ok := moveCrossCursor(p.cursorCross, direction)
	if !ok {
		return
	}
	p.cursorCross = dest
}

type phaseInitialRoad struct {
	game        *Game
	playerTurn  int
	sourceCross board.CrossCoord
	cursorCross board.CrossCoord
}

func PhaseInitialRoad(game *Game, playerTurn int, sourceCross board.CrossCoord) Phase {
	return &phaseInitialRoad{
		game:        game,
		playerTurn:  playerTurn,
		sourceCross: sourceCross,
		cursorCross: sourceCross.Neighbors()[0],
	}
}

func (p *phaseInitialRoad) Confirm() Phase {
	roadCoord := board.NewPathCoord(p.sourceCross, p.cursorCross)
	p.game.board.SetRoad(roadCoord, p.game.players[p.playerTurn].Name)
	// TODO: transition to next phase
	return p
}

func (p *phaseInitialRoad) Cancel() Phase {
	return PhaseInitialSettlements(p.game, p.playerTurn)
}

func (p *phaseInitialRoad) HelpText() string {
	return fmt.Sprintf("%s's turn. Place your initial road connected to the settlement with 'enter'.", p.game.players[p.playerTurn].Name)
}

func (p *phaseInitialRoad) CurrentCursor() interface{} {
	return p.cursorCross
}

func (p *phaseInitialRoad) MoveCursor(direction string) {
	var dest board.CrossCoord
	var ok bool

	switch direction {
	case "up":
		dest, ok = p.sourceCross.Up()
	case "down":
		dest, ok = p.sourceCross.Down()
	case "left":
		dest, ok = p.sourceCross.Left()
	case "right":
		dest, ok = p.sourceCross.Right()
	default:
		return
	}

	if ok {
		p.cursorCross = dest
	}
}
