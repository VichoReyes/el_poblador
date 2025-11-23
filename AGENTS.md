# El Poblador - Project Context

## Project Overview
El Poblador is a Go implementation of the Catan board game, featuring a hexagonal board, turn-based gameplay, and various game phases.

## Tasks
Things to be done can be found in TASKS.md. As you complete tasks, remove them from this file. Emphasis on _remove_, not "put in a completed section".
If the completion of a task uncoverered a new task to be added, add it. You may put some small guidelines to provide continuity.

## Project Structure
- **`board/`** - Board representation, coordinates, rendering, and terrain
  - `board.go` - Core Board struct with tiles, roads, settlements, cities, and robber
  - `coordinates.go` - TileCoord, CrossCoord, PathCoord coordinate systems
  - `terrain.go` - Tile types and resource generation
  - `render.go` - Board visualization
  - `robber_*.go` - Robber placement and blocking mechanics
- **`game/`** - Core game logic, phases, players, and turn management  
  - `game.go` - Main Game struct and game loop
  - `turn.go` - All game phases (dice roll, building, robber placement, etc.)
  - `player.go` - Player resources, development cards, and actions
  - `dev_card.go` - Development card types and effects
  - `initial_settlements.go` - Game setup phase
- **`cmd/`** - Command-line applications for testing

## Key Files
- **`.build.yml`** - Build manifest containing commands to verify the project works
- **`main.go`** - Main entry point
- **`go.mod`** - Go module dependencies

## Verification Commands
To verify if things work after a change, run the commands in `.build.yml`:

```bash
cd el_poblador
go test -v ./...
go run cmd/testboard/main.go
go run cmd/fullscreen/main.go
```

**Note**: The commands in the build manifest can change, so always refer to the current `.build.yml` file for the latest verification steps.

## Documentation
- Each package may contain markdown files with additional documentation
- Check package directories for `.md` files that provide context-specific information
- Example: `board/coordinates.md` explains the coordinate system and rendering approach

## Development Workflow
1. Make changes to the codebase
2. Run the verification commands from `.build.yml`
3. Check for any markdown files in relevant package directories for additional context
4. Ensure tests pass and applications run correctly

## Development Guidelines
- **Interface Design**: Avoid adding methods to interfaces if they won't really be necessary for most implementations
- **Simplicity**: Prefer simple, static implementations over complex dynamic ones to reduce bugs
- **Testing**: Write tests for new functionality to ensure reliability

## Testing Strategy
- Tests are located alongside source files with `_test.go` suffix
- Use `go test -v ./...` to run all tests
- Test applications (`cmd/testboard/main.go`, `cmd/fullscreen/main.go`) help verify functionality

## Architecture Patterns

### Game Phase System
The game uses a phase-based architecture where all game states implement the `Phase` interface:
- `Phase` - Base interface with `Confirm()`, `MoveCursor()`, `BoardCursor()`, `HelpText()`
- `PhaseWithMenu` - Phases that show selection menus (extends Phase)
- `PhaseCancelable` - Phases that can be cancelled (extends Phase)

Key phases include:
- `phaseDiceRoll` - Start of turn (roll dice or play Knight)
- `phasePlaceRobber` - Place robber after rolling 7 or playing Knight card
- `phaseStealCard` - Select player to steal from after robber placement
- `phaseIdle` - Main turn menu (build, trade, play dev cards, end turn)
- `phaseBuilding` - Building selection submenu
- Various placement phases (`phaseRoadStart`, `phaseSettlementPlacement`, etc.)

### Coordinate System
- **TileCoord** - Hexagonal tiles, odd sum (x+y) coordinates
- **CrossCoord** - Intersection points where settlements/cities are placed
- **PathCoord** - Edges between crosses where roads are placed

### Robber Mechanics
- Stored as `TileCoord` in `Board.robber`
- `PlaceRobber(coord)` moves robber and returns adjacent stealable player IDs
- `GetRobber()` returns current robber position
- **Rule**: Cannot place robber on the same tile it's already on (enforced in `phasePlaceRobber.Confirm()`)

### Resource Generation
- `Board.GenerateResources(sum)` returns `map[int][]ResourceType` (player ID -> resources)
- Tiles with robber are skipped during resource generation
- Settlements get 1 resource, cities get 2 resources per matching tile

### Development Cards
- Stored in `Player.hiddenDevCards` (playable) and `Player.playedDevCards`
- Knight cards trigger robber placement phase
- Road Building gives 2 free roads
- Monopoly and Year of Plenty have dedicated phases

## Package-Specific Context
When working on specific packages, always check for markdown files in the package directory that may contain:
- Design decisions and rationale
- Implementation details
- Coordinate systems and algorithms
- Rendering approaches
- Testing strategies
- avoid committing compiled binaries. before making a commit, check whether you need to add any to gitignore