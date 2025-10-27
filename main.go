package main

import (
	"bytes"
	"el_poblador/game"
	"encoding/gob"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	game           *game.Game
	width          int
	height         int
	userPlayer     *int
	twoColumnCycle int // 0-1: for width 90-119
	oneColumnCycle int // 0-2: for width <90
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
		case "tab":
			m.twoColumnCycle = (m.twoColumnCycle + 1) % 2
			m.oneColumnCycle = (m.oneColumnCycle + 1) % 3
		case "up", "down", "left", "right":
			m.game.MoveCursor(msg.String(), m.userPlayer)
		case "enter":
			m.game.ConfirmAction(m.userPlayer)
			if m.game.ShouldQuit() {
				return m, tea.Quit
			}
		case "esc":
			m.game.CancelAction(m.userPlayer)
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

func (m model) View() string {
	return m.game.Print(m.width, m.height, m.userPlayer, m.twoColumnCycle, m.oneColumnCycle)
}

func loadGameState(filename string) (*game.Game, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read failed: %w", err)
	}

	var g game.Game
	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&g); err != nil {
		return nil, fmt.Errorf("decoding failed: %w", err)
	}

	// Restore phase to PhaseDiceRoll
	g.SetPhase(game.PhaseDiceRoll(&g))

	return &g, nil
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  el_poblador new <player1> <player2> <player3> [player4]")
	fmt.Println("  el_poblador load <filename.gob>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  new   Start a new game with 3-4 players")
	fmt.Println("  load  Load a saved game from file")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	var g *game.Game

	switch command {
	case "new":
		if len(os.Args) < 5 || len(os.Args) > 6 {
			fmt.Println("Error: 'new' command requires 3-4 player names")
			fmt.Println()
			printUsage()
			os.Exit(1)
		}
		g = &game.Game{}
		g.Start(os.Args[2:])

	case "load":
		if len(os.Args) != 3 {
			fmt.Println("Error: 'load' command requires a filename")
			fmt.Println()
			printUsage()
			os.Exit(1)
		}
		loadedGame, err := loadGameState(os.Args[2])
		if err != nil {
			fmt.Printf("Failed to load game: %v\n", err)
			os.Exit(1)
		}
		g = loadedGame

	default:
		fmt.Printf("Error: unknown command '%s'\n", command)
		fmt.Println()
		printUsage()
		os.Exit(1)
	}

	p := tea.NewProgram(model{game: g}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
