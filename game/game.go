package game

import (
	"el_poblador/board"
	"math/rand/v2"
)

type Player struct {
	Name  string
	color int // 8 bit color code
}

type GamePhase int

const (
	GameSetup GamePhase = iota
	GameInitialSettlements
	GamePlaying
	GameEnded
)

type TurnPhase int

const (
	BeforeDice TurnPhase = iota
	Robbing
	TurnMain
)

// Meaningful zero value
type Game struct {
	board      *board.Board
	players    []Player
	playerTurn int
	gamePhase  GamePhase
	turnPhase  TurnPhase
	lastDice   [2]int
}

func (g *Game) Start(playerNames []string) {
	if len(playerNames) < 3 || len(playerNames) > 4 {
		panic("Game must have 3-4 players")
	}
	if g.gamePhase != GameSetup {
		panic("Cannot start a game that has already started")
	}
	colors := []int{20, 88, 165, 103} // blue, red, purple, white
	g.players = make([]Player, len(playerNames))
	for i, name := range playerNames {
		g.players[i] = Player{Name: name, color: colors[i]}
	}
	rand.Shuffle(len(g.players), func(i, j int) {
		g.players[i], g.players[j] = g.players[j], g.players[i]
	})
	g.board = board.NewLegalBoard()
	g.gamePhase = GameInitialSettlements
}
