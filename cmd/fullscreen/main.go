package main

import (
	"el_poblador/board"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// model represents the state of our application.
type model struct {
	board  *board.Board
	width  int
	height int
}

// initialModel creates the initial state of the application.
func initialModel() model {
	return model{
		board: board.NewChaoticBoard(),
	}
}

// Init is the first function that will be called. It can be used to send
// an initial command. We don't need to do that here, so we'll return nil.
func (m model) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received. Use it to update your model.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "c" key creates a chaotic board.
		case "c":
			m.board = board.NewChaoticBoard()
			return m, nil

		}

	// Is it a window resize?
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View is called to render the program's UI.
func (m model) View() string {
	margin := lipgloss.NewStyle().Margin(1)
	// Title text styled with bold.
	titleText := lipgloss.NewStyle().Bold(true).Render("El Poblador - Bubble Tea Version")

	// The title is placed in the center of the screen, with a margin below it.
	renderedTitle := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, titleText)

	// Help text styled as faint.
	helpText := lipgloss.NewStyle().Faint(true).Render("Press 'c' or 'l' to regenerate the board. Press 'q' to quit.")

	// The help text is placed in the center of the screen.
	renderedHelp := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, helpText)

	// Board
	boardLines := m.board.Print(nil)
	boardContent := strings.Join(boardLines, "\n")

	// players
	players := strings.Join([]string{"John", "Jane", "Jim (turn)", "Jill"}, "\n\n")
	dice := "⚂ ⚄"
	sidebar := margin.Render(lipgloss.JoinVertical(lipgloss.Left, dice, players))
	renderedPlayers := lipgloss.JoinHorizontal(lipgloss.Top, boardContent, sidebar)

	// Calculate the available space for the board.
	availableHeight := m.height - lipgloss.Height(renderedTitle) - lipgloss.Height(renderedHelp)

	// Main content is the board, centered in the available space.
	mainContent := lipgloss.Place(m.width, availableHeight,
		lipgloss.Center,
		lipgloss.Center,
		renderedPlayers,
	)

	return lipgloss.JoinVertical(lipgloss.Left, renderedTitle, mainContent, renderedHelp)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
