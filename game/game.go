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

// Meaningful zero value
type Game struct {
	board      *board.Board
	players    []Player
	playerTurn int
	gamePhase  GamePhase
	turnPhase  TurnPhase
	lastDice   [2]int
}

func (g *Game) Print(width, height int) string {
	margin := lipgloss.NewStyle().Margin(1)
	// Title text styled with bold.
	titleText := lipgloss.NewStyle().Bold(true).Render("El Poblador - Bubble Tea Version")

	// The title is placed in the center of the screen, with a margin below it.
	renderedTitle := lipgloss.PlaceHorizontal(width, lipgloss.Center, titleText)

	// Help text styled as faint.
	helpText := lipgloss.NewStyle().Faint(true).Render("Press 'c' or 'l' to regenerate the board. Press 'q' to quit.")

	// The help text is placed in the center of the screen.
	renderedHelp := lipgloss.PlaceHorizontal(width, lipgloss.Center, helpText)

	// Board
	boardLines := g.board.Print()
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
	availableHeight := height - lipgloss.Height(renderedTitle) - lipgloss.Height(renderedHelp)

	// Main content is the board, centered in the available space.
	mainContent := lipgloss.Place(width, availableHeight,
		lipgloss.Center,
		lipgloss.Center,
		renderedPlayers,
	)

	return lipgloss.JoinVertical(lipgloss.Left, renderedTitle, mainContent, renderedHelp)
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
