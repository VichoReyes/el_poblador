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
	board      *board.Board
	players    []Player
	lastDice   [2]int
	phase      Phase
	playerTurn int
}

// requestPlayer is the player that the user is playing as.
// If nil, the game will render from the perspective of the turn holder.
func (g *Game) Print(width, height int, requestPlayer *int) string {
	playerPerspective := g.playerPerspective(requestPlayer)
	margin := lipgloss.NewStyle().Margin(1)

	help := g.helpText(width)

	boardLines := g.board.Print(g.phase.BoardCursor())
	boardContent := strings.Join(boardLines, "\n")

	var dice string
	if g.lastDice[0] != 0 {
		dice = fmt.Sprintf("Dice: %d (%d + %d)", g.lastDice[0]+g.lastDice[1], g.lastDice[0], g.lastDice[1])
		dice = margin.Render(dice)
	}

	var playerList []string
	for i, player := range g.players {
		var name string
		if i == playerPerspective {
			name = fmt.Sprintf("%s (you)", player.Render(player.Name))
		} else {
			name = player.Render(player.Name)
		}
		info := player.Render(fmt.Sprintf(" has %d resources", player.TotalResources()))
		playerList = append(playerList, name, info)
	}
	otherPlayers := margin.Render(strings.Join(playerList, "\n"))

	myPlayer := g.players[playerPerspective]
	myResources := []string{"Your resources:"}
	for resource, amount := range myPlayer.resources {
		myResources = append(myResources, fmt.Sprintf("%s: %d", resource, amount))
	}
	myResourcesStr := margin.Render(strings.Join(myResources, "\n"))

	var phaseSidebar string
	if p, ok := g.phase.(PhaseWithMenu); ok {
		if playerPerspective == g.playerTurn {
			phaseSidebar = margin.Render(p.Menu())
		}
	}
	sidebar := lipgloss.JoinVertical(lipgloss.Left, dice, otherPlayers, myResourcesStr, phaseSidebar)
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

func (g *Game) playerPerspective(requestPlayer *int) int {
	if requestPlayer != nil && *requestPlayer < len(g.players) {
		return *requestPlayer
	} else {
		return g.playerTurn
	}
}

func (g *Game) helpText(width int) string {
	player := g.players[g.playerTurn]
	help := fmt.Sprintf("%s's turn. %s", player.Render(player.Name), g.phase.HelpText())
	renderedHelp := lipgloss.PlaceHorizontal(width, lipgloss.Center, help)
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
		return g.players[playerId].Render(content)
	})
	g.playerTurn = 0
	g.phase = PhaseInitialSettlements(g, true)
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

func (g *Game) MoveCursor(direction string, requestPlayer *int) {
	playerPerspective := g.playerPerspective(requestPlayer)
	// for now simply ignore moves from other players
	if playerPerspective != g.playerTurn {
		return
	}
	g.phase.MoveCursor(direction)
}

// Function for testing purposes: move the cursor to any valid settlement location
func (g *Game) MoveCursorToPlaceSettlement() {
	find := func() board.CrossCoord {
		for x := 0; x <= 5; x++ {
			for y := 0; y <= 10; y++ {
				coord, valid := board.NewCrossCoord(x, y)
				if valid && g.board.CanPlaceSettlement(coord) {
					return coord
				}
			}
		}
		panic("no valid settlement location found")
	}
	if p, ok := g.phase.(*phaseInitialSettlements); ok {
		p.cursorCross = find()
	}
}

func (g *Game) ConfirmAction(requestPlayer *int) {
	playerPerspective := g.playerPerspective(requestPlayer)
	// for now simply ignore actions from other players
	if playerPerspective != g.playerTurn {
		return
	}
	g.phase = g.phase.Confirm()
}

func (g *Game) CancelAction(requestPlayer *int) {
	playerPerspective := g.playerPerspective(requestPlayer)
	// for now simply ignore actions from other players
	if playerPerspective != g.playerTurn {
		return
	}
	g.phase = g.phase.Cancel()
}
