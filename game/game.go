package game

import (
	"el_poblador/board"
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/charmbracelet/lipgloss"
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
	OfferingTrade
)

type CursorTarget int

const (
	TargetCross CursorTarget = iota
	TargetTile
	TargetCard
)

// Meaningful zero value
type Game struct {
	board      *board.Board
	players    []Player
	playerTurn int
	winner     *Player
	gamePhase  GamePhase
	turnPhase  TurnPhase
	lastDice   [2]int

	cursorTarget CursorTarget
	cursorCross  board.CrossCoord
	cursorTile   board.TileCoord
}

func (g *Game) Print(width, height int) string {
	margin := lipgloss.NewStyle().Margin(1)

	help := g.helpText(width)

	// Board
	var boardLines []string
	switch g.cursorTarget {
	case TargetCross:
		boardLines = g.board.Print(g.cursorCross)
	case TargetTile:
		boardLines = g.board.Print(g.cursorTile)
	default:
		boardLines = g.board.Print(nil)
	}
	boardContent := strings.Join(boardLines, "\n")

	// players
	playerNames := make([]string, len(g.players))
	for i, player := range g.players {
		playerNames[i] = fmt.Sprintf("\033[38;5;%dm %s \033[0m", player.color, player.Name)
	}
	players := strings.Join(playerNames, "\n\n")
	dice := "⚂ ⚄"
	sidebar := margin.Render(lipgloss.JoinVertical(lipgloss.Left, dice, players))
	renderedPlayers := lipgloss.JoinHorizontal(lipgloss.Top, boardContent, sidebar)

	// Calculate the available space for the board.
	availableHeight := height - lipgloss.Height(help)

	// Main content is the board, centered in the available space.
	mainContent := lipgloss.Place(width, availableHeight,
		lipgloss.Center,
		lipgloss.Center,
		renderedPlayers,
	)

	return lipgloss.JoinVertical(lipgloss.Left, mainContent, help)
}

func (g *Game) helpText(width int) string {
	var text string
	switch g.gamePhase {
	case GameSetup:
		text = "You shouldn't see this"
	case GameInitialSettlements:
		text = fmt.Sprintf(
			"%s's turn. Place your initial settlements and roads on the board with 'enter'.",
			g.players[g.playerTurn].Name,
		)
	case GamePlaying:
		text = fmt.Sprintf("%s's turn. ", g.players[g.playerTurn].Name)
		switch g.turnPhase {
		case BeforeDice:
			text += "Press 'd' to roll the dice."
		case Robbing:
			text += "Pick a tile to place the robber on."
		case TurnMain:
			text += "Press 'tab' to switch between the main actions, or 'n' to end your turn."
		default:
			panic("unknown turn phase")
		}
	case GameEnded:
		text = fmt.Sprintf("Finished! %s wins! Press 'q' to quit.", g.winner.Name)
	default:
		panic("unknown game phase")
	}
	style := lipgloss.NewStyle().Faint(true)
	renderedHelp := lipgloss.PlaceHorizontal(width, lipgloss.Center, style.Render(text))
	return renderedHelp
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
	g.cursorTarget = TargetCross
	g.cursorCross = g.board.ValidCrossCoord()
}

func (g *Game) MoveCursor(direction string) {
	var dest board.CrossCoord
	var ok bool
	switch direction {
	case "up":
		dest, ok = g.cursorCross.Up()
	case "down":
		dest, ok = g.cursorCross.Down()
	case "left":
		dest, ok = g.cursorCross.Left()
	case "right":
		dest, ok = g.cursorCross.Right()
	}
	if ok {
		g.cursorCross = dest
	}
}
