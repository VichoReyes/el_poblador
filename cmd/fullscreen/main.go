package main

import (
	"bufio"
	"el_poblador/board"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// Terminal represents our full-screen terminal application
type Terminal struct {
	originalState *term.State
	width         int
	height        int
	board         *board.Board
}

// NewTerminal creates a new terminal instance
func NewTerminal() (*Terminal, error) {
	// Get terminal size
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return nil, fmt.Errorf("failed to get terminal size: %w", err)
	}

	// Save original terminal state
	originalState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return nil, fmt.Errorf("failed to enter raw mode: %w", err)
	}

	return &Terminal{
		originalState: originalState,
		width:         width,
		height:        height,
		board:         board.NewChaoticBoard(),
	}, nil
}

// Cleanup restores the terminal to its original state
func (t *Terminal) Cleanup() {
	// Show cursor
	fmt.Print("\033[?25h")
	// Exit alternate screen buffer
	fmt.Print("\033[?1049l")
	// Restore original terminal state
	term.Restore(int(os.Stdin.Fd()), t.originalState)
}

// Clear clears the screen and moves cursor to top-left
func (t *Terminal) Clear() {
	fmt.Print("\033[2J\033[H")
}

// HideCursor hides the terminal cursor
func (t *Terminal) HideCursor() {
	fmt.Print("\033[?25l")
}

// ShowCursor shows the terminal cursor
func (t *Terminal) ShowCursor() {
	fmt.Print("\033[?25h")
}

// MoveCursor moves the cursor to the specified position (1-indexed)
func (t *Terminal) MoveCursor(row, col int) {
	fmt.Printf("\033[%d;%dH", row, col)
}

// EnterAlternateScreen enters the alternate screen buffer
func (t *Terminal) EnterAlternateScreen() {
	fmt.Print("\033[?1049h")
}

// DrawBoard draws the game board centered on screen
func (t *Terminal) DrawBoard() {
	lines := t.board.Print()

	// Calculate centering offset
	maxLineLength := 0
	for _, line := range lines {
		if len(line) > maxLineLength {
			maxLineLength = len(line)
		}
	}

	startRow := max((t.height-len(lines))/2, 1)

	startCol := max((t.width-maxLineLength)/2, 1)

	// Draw each line
	for i, line := range lines {
		t.MoveCursor(startRow+i, startCol)
		fmt.Print(line)
	}
}

// DrawStatusBar draws a status bar at the bottom of the screen
func (t *Terminal) DrawStatusBar() {
	t.MoveCursor(t.height, 1)
	statusMsg := fmt.Sprintf("Terminal: %dx%d | Press 'q' to quit, 'r' to regenerate board", t.width, t.height)

	// Pad or truncate to fit screen width
	if len(statusMsg) > t.width {
		statusMsg = statusMsg[:t.width]
	} else {
		statusMsg = statusMsg + strings.Repeat(" ", t.width-len(statusMsg))
	}

	// Draw with reverse video (white on black)
	fmt.Printf("\033[7m%s\033[0m", statusMsg)
}

// DrawTitle draws a title at the top of the screen
func (t *Terminal) DrawTitle() {
	title := "El Poblador - Settlers of Catan"
	startCol := (t.width - len(title)) / 2
	if startCol < 1 {
		startCol = 1
	}

	t.MoveCursor(1, startCol)
	fmt.Printf("\033[1m%s\033[0m", title) // Bold text
}

// Render redraws the entire screen
func (t *Terminal) Render() {
	t.Clear()
	t.DrawTitle()
	t.DrawBoard()
	t.DrawStatusBar()
}

// ReadKey reads a single keypress
func (t *Terminal) ReadKey() (byte, error) {
	reader := bufio.NewReader(os.Stdin)
	char, err := reader.ReadByte()
	return char, err
}

// Run starts the main application loop
func (t *Terminal) Run() error {
	// Enter alternate screen and hide cursor
	t.EnterAlternateScreen()
	t.HideCursor()

	// Setup signal handling for cleanup
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Initial render
	t.Render()

	// Main event loop
	for {
		select {
		case <-sigChan:
			return nil
		default:
			// Read key input
			key, err := t.ReadKey()
			if err != nil {
				return fmt.Errorf("failed to read key: %w", err)
			}

			// Handle input
			switch key {
			case 'q', 'Q':
				return nil
			case 'r', 'R':
				t.board = board.NewChaoticBoard()
				t.Render()
			case '\x03': // Ctrl+C
				return nil
			}
		}
	}
}

func main() {
	// Create terminal instance
	terminal, err := NewTerminal()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating terminal: %v\n", err)
		os.Exit(1)
	}

	// Ensure cleanup happens
	defer terminal.Cleanup()

	// Run the application
	if err := terminal.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running terminal: %v\n", err)
		os.Exit(1)
	}
}
