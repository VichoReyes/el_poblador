# El Poblador - Settlers of Catan TUI

A terminal-based implementation of Settlers of Catan in Go.

![game screenshot](screenshot.png)

## Running the Project

### Main Game

```bash
go run main.go <player1> <player2> <player3> [player4]
```

Runs the main Catan game with 3-4 players. Provide player names as command-line arguments.

**Controls:**
- Arrow keys: Move cursor
- Enter: Confirm action
- Esc: Cancel action (not always available)
- 1-4: Switch to specific player's perspective
- 0: Switch back to current turn holder's perspective
- q/Ctrl+C: Quit game

### Preview Board
```bash
go run cmd/testboard/main.go
```
Displays a test board with sample game state.

### Full-Screen Terminal
```bash
go run cmd/fullscreen/main.go
```
Shows the whole game interface. This is how the screenshot above was generated.

### Running Tests
```bash
go test -v ./...
```

Runs all tests across the project.

## Project Structure

- **`board/`** - Board representation, coordinates, rendering, and terrain
- **`game/`** - Core game logic, phases, players, and turn management
- **`cmd/`** - Command-line applications for testing

## Game Features

- Full Catan gameplay with dice rolling, building, trading, and development cards
- Hexagonal board with proper coordinate system
- Turn-based phase system (dice roll, building, robber placement, etc.)
- Resource generation and management
- Victory point tracking
- Robber mechanics with card stealing
- Development cards (Knight, Road Building, Monopoly, Year of Plenty, Victory Points)

## Missing

(see Issues for a more up to date list)

- Longest Road and Largest Army
- Trade
  - Trading ports
- Online Multiplayer
- Better localization (eliminating spanglish)
- 5 to 6 player mode

## License

EUPL v1.2, in spanish