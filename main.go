package main

import (
	"el_poblador/game"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	game   game.Game
	width  int
	height int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "down", "left", "right":
			m.game.MoveCursor(msg.String())
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	}
	return m, nil
}

func (m model) View() string {
	return m.game.Print(m.width, m.height)
}

func main() {
	game := game.Game{}
	if len(os.Args) < 4 || len(os.Args) > 5 {
		fmt.Println("Please provide 3-4 player names as arguments")
		os.Exit(1)
	}
	game.Start(os.Args[1:])
	p := tea.NewProgram(model{game: game}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
