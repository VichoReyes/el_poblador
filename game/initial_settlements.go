package game

import (
	"el_poblador/board"
	"fmt"
)

func nextInitialPhase(game *Game, currentPlayer int, isFirstPair bool) Phase {
	numPlayers := len(game.players)

	if isFirstPair {
		nextPlayer := currentPlayer + 1
		if nextPlayer < numPlayers {
			return PhaseInitialSettlements(game, nextPlayer, true)
		} else {
			return PhaseInitialSettlements(game, numPlayers-1, false)
		}
	} else {
		nextPlayer := currentPlayer - 1
		if nextPlayer >= 0 {
			return PhaseInitialSettlements(game, nextPlayer, false)
		} else {
			return PhaseDiceRoll(game, 0)
		}
	}
}

type phaseInitialSettlements struct {
	game        *Game
	playerTurn  int
	cursorCross board.CrossCoord
	isFirstPair bool
}

func PhaseInitialSettlements(game *Game, playerTurn int, isFirstPair bool) Phase {
	cursorCross := game.board.ValidCrossCoord()
	return &phaseInitialSettlements{game: game, playerTurn: playerTurn, cursorCross: cursorCross, isFirstPair: isFirstPair}
}

func (p *phaseInitialSettlements) Confirm() Phase {
	if !p.game.board.SetSettlement(p.cursorCross, p.playerTurn) {
		return p
	}
	if !p.isFirstPair {
		player := p.game.players[p.playerTurn]
		adjacentTiles := p.game.board.AdjacentTiles(p.cursorCross)
		for _, tile := range adjacentTiles {
			resource, ok := board.TileResource(tile)
			if ok {
				player.AddResource(resource)
			}
		}
	}
	return PhaseInitialRoad(p.game, p.playerTurn, p.cursorCross, p.isFirstPair)
}

func (p *phaseInitialSettlements) Cancel() Phase {
	return p
}

func (p *phaseInitialSettlements) HelpText() string {
	var num string
	if p.isFirstPair {
		num = "first"
	} else {
		num = "second"
	}
	return fmt.Sprintf("%s's turn. Place your %s settlement on the board with 'enter'.",
		p.game.players[p.playerTurn].Name,
		num,
	)
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
	playerTurn  int
	sourceCross board.CrossCoord
	cursorCross board.CrossCoord
	isFirstPair bool
}

func PhaseInitialRoad(game *Game, playerTurn int, sourceCross board.CrossCoord, isFirstPair bool) Phase {
	return &phaseInitialRoad{
		game:        game,
		playerTurn:  playerTurn,
		sourceCross: sourceCross,
		cursorCross: sourceCross.Neighbors()[0],
		isFirstPair: isFirstPair,
	}
}

func (p *phaseInitialRoad) Confirm() Phase {
	roadCoord := board.NewPathCoord(p.sourceCross, p.cursorCross)
	p.game.board.SetRoad(roadCoord, p.playerTurn)
	return nextInitialPhase(p.game, p.playerTurn, p.isFirstPair)
}

func (p *phaseInitialRoad) Cancel() Phase {
	return PhaseInitialSettlements(p.game, p.playerTurn, p.isFirstPair)
}

func (p *phaseInitialRoad) HelpText() string {
	return fmt.Sprintf("%s's turn. Place your initial road connected to the settlement with 'enter'.", p.game.players[p.playerTurn].Name)
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
