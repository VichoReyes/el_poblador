package main

import (
	"bytes"
	"el_poblador/game"
	"encoding/gob"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	game        *game.Game
	width       int
	height      int
	userPlayer  *int
	saveMessage string
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
			m.game.MoveCursor(msg.String(), m.userPlayer)
		case "enter":
			m.game.ConfirmAction(m.userPlayer)
		case "esc":
			m.game.CancelAction(m.userPlayer)
		case "s":
			// Save game state
			if err := saveGameState(m.game); err != nil {
				m.saveMessage = fmt.Sprintf("Save failed: %v", err)
			} else {
				m.saveMessage = "Game saved successfully!"
			}
		// switch to specific player's perspective
		case "1", "2", "3", "4":
			player := int(msg.String()[0] - '1')
			m.userPlayer = &player
		// switch back to turn holder's perspective
		case "0":
			m.userPlayer = nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	}
	return m, nil
}

func saveGameState(g *game.Game) error {
	// Create filename with timestamp
	filename := fmt.Sprintf("game_save_%s.gob", time.Now().Format("2006-01-02_15-04-05"))

	// Encode game to gob
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(g); err != nil {
		return fmt.Errorf("encoding failed: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filename, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("write failed: %w", err)
	}

	return nil
}

func (m model) View() string {
	view := m.game.Print(m.width, m.height, m.userPlayer)
	if m.saveMessage != "" {
		view = view + "\n" + m.saveMessage
	}
	return view
}

func main() {
	game := game.Game{}
	if len(os.Args) < 4 || len(os.Args) > 5 {
		fmt.Println("Please provide 3-4 player names as arguments")
		os.Exit(1)
	}
	game.Start(os.Args[1:])
	p := tea.NewProgram(model{game: &game}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
