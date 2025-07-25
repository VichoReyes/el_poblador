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
