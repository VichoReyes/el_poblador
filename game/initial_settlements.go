package game

import (
	"el_poblador/board"
	"fmt"
)

func nextInitialPhase(game *Game, isFirstPair bool) Phase {
	numPlayers := len(game.players)

	if isFirstPair {
		if game.playerTurn == numPlayers-1 {
			return PhaseInitialSettlements(game, false)
		} else {
			game.playerTurn++
			return PhaseInitialSettlements(game, true)
		}
	} else {
		if game.playerTurn == 0 {
			return PhaseDiceRoll(game)
		} else {
			game.playerTurn--
			return PhaseInitialSettlements(game, false)
		}
	}
}

type phaseInitialSettlements struct {
	game        *Game
	cursorCross board.CrossCoord
	isFirstPair bool
}

func PhaseInitialSettlements(game *Game, isFirstPair bool) Phase {
	cursorCross := game.board.ValidCrossCoord()
	return &phaseInitialSettlements{game: game, cursorCross: cursorCross, isFirstPair: isFirstPair}
}

func (p *phaseInitialSettlements) Confirm() Phase {
	if !p.game.board.SetSettlement(p.cursorCross, p.game.playerTurn) {
		return p
	}
	if !p.isFirstPair {
		player := p.game.players[p.game.playerTurn]
		adjacentTiles := p.game.board.AdjacentTiles(p.cursorCross)
		for _, tile := range adjacentTiles {
			resource, ok := board.TileResource(tile)
			if ok {
				player.AddResource(resource)
			}
		}
	}
	
	// Check for game end after building initial settlement (unlikely but for completeness)
	if winner := p.game.CheckGameEnd(); winner != nil {
		return PhaseGameEnd(p.game, winner)
	}
	
	return PhaseInitialRoad(p.game, p.cursorCross, p.isFirstPair)
}

func (p *phaseInitialSettlements) HelpText() string {
	var num string
	if p.isFirstPair {
		num = "first"
	} else {
		num = "second"
	}
	return fmt.Sprintf("Place your %s settlement on the board with 'enter'.", num)
}

func (p *phaseInitialSettlements) BoardCursor() interface{} {
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
	sourceCross board.CrossCoord
	cursorCross board.CrossCoord
	isFirstPair bool
}

func PhaseInitialRoad(game *Game, sourceCross board.CrossCoord, isFirstPair bool) Phase {
	return &phaseInitialRoad{
		game:        game,
		sourceCross: sourceCross,
		cursorCross: sourceCross.Neighbors()[0],
		isFirstPair: isFirstPair,
	}
}

func (p *phaseInitialRoad) Confirm() Phase {
	roadCoord := board.NewPathCoord(p.sourceCross, p.cursorCross)
	p.game.board.SetRoad(roadCoord, p.game.playerTurn)
	return nextInitialPhase(p.game, p.isFirstPair)
}

func (p *phaseInitialRoad) Cancel() Phase {
	return PhaseInitialSettlements(p.game, p.isFirstPair)
}

func (p *phaseInitialRoad) HelpText() string {
	return "Place a road connected to the settlement by selecting its direction."
}

func (p *phaseInitialRoad) BoardCursor() interface{} {
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
