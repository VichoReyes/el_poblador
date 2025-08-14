# El Poblador - Project Context

## Project Overview
El Poblador is a Go implementation of the Catan board game, featuring a hexagonal board, turn-based gameplay, and various game phases.

## Project Structure
- **`board/`** - Board representation, coordinates, rendering, and terrain
- **`game/`** - Core game logic, phases, players, and turn management
- **`cmd/`** - Command-line applications for testing and fullscreen display

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

## Testing Strategy
- Tests are located alongside source files with `_test.go` suffix
- Use `go test -v ./...` to run all tests
- Test applications (`cmd/testboard/main.go`, `cmd/fullscreen/main.go`) help verify functionality

## Package-Specific Context
When working on specific packages, always check for markdown files in the package directory that may contain:
- Design decisions and rationale
- Implementation details
- Coordinate systems and algorithms
- Rendering approaches
- Testing strategies
