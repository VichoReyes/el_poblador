package game

import (
	"el_poblador/board"
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Phase interface {
	Confirm() Phase
	Cancel() Phase
	MoveCursor(direction string)
	BoardCursor() interface{}
	HelpText() string
}

type PhaseWithMenu interface {
	Phase
	Menu() string
}

type Game struct {
	board    *board.Board
	players  []Player
	lastDice [2]int
	phase    Phase
}

func (g *Game) Print(width, height int) string {
	margin := lipgloss.NewStyle().Margin(1)

	help := g.helpText(width)

	// Board
	boardLines := g.board.Print(g.phase.BoardCursor())
	boardContent := strings.Join(boardLines, "\n")

	// players
	playerNames := make([]string, len(g.players))
	for i, player := range g.players {
		playerNames[i] = fmt.Sprintf("\033[38;5;%dm %s \033[0m", player.color, player.Name)
	}
	players := strings.Join(playerNames, "\n\n")
	dice := "⚂ ⚄"
	var phaseSidebar string
	if p, ok := g.phase.(PhaseWithMenu); ok {
		phaseSidebar = p.Menu()
	}
	sidebar := margin.Render(lipgloss.JoinVertical(lipgloss.Left, dice, players, phaseSidebar))
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
	text := g.phase.HelpText()
	style := lipgloss.NewStyle().Faint(true)
	renderedHelp := lipgloss.PlaceHorizontal(width, lipgloss.Center, style.Render(text))
	return renderedHelp
}

func (g *Game) Start(playerNames []string) {
	if len(playerNames) < 3 || len(playerNames) > 4 {
		panic("Game must have 3-4 players")
	}
	colors := []int{20, 88, 165, 103} // blue, red, purple, white
	g.players = make([]Player, len(playerNames))
	for i, name := range playerNames {
		g.players[i] = Player{Name: name, color: colors[i], resources: make(map[board.ResourceType]int)}
	}
	rand.Shuffle(len(g.players), func(i, j int) {
		g.players[i], g.players[j] = g.players[j], g.players[i]
	})
	g.board = board.NewLegalBoard(func(playerId int, content string) string {
		return fmt.Sprintf("\033[38;5;%dm%s\033[0m", g.players[playerId].color, content)
	})
	g.phase = PhaseInitialSettlements(g, 0, true)
}

func moveCrossCursor(from board.CrossCoord, direction string) (dest board.CrossCoord, ok bool) {
	switch direction {
	case "up":
		dest, ok = from.Up()
	case "down":
		dest, ok = from.Down()
	case "left":
		dest, ok = from.Left()
	case "right":
		dest, ok = from.Right()
	default:
		panic("unknown direction")
	}
	return dest, ok
}

func (g *Game) MoveCursor(direction string) {
	g.phase.MoveCursor(direction)
}

func (g *Game) ConfirmAction() {
	g.phase = g.phase.Confirm()
}

func (g *Game) CancelAction() {
	g.phase = g.phase.Cancel()
}
