# El Poblador - Settlers of Catan TUI

A terminal-based implementation of Settlers of Catan in Go.

## Project Structure

```
el_poblador/
├── main.go                 # Main game executable
├── board/
│   └── board.go           # Board component with hexagon printing
├── cmd/
│   └── testboard/
│       └── main.go        # Board component test executable
├── go.mod                 # Go module file
└── README.md              # This file
```

## Running the Project

### Main Game

```bash
go run main.go
```

Runs the main game executable that imports and uses the board functionality.

### Test Board Component
```bash
go run cmd/testboard/main.go
```
Runs the board component in isolation for testing - this is Go's equivalent to Python's `if __name__ == "__main__"` pattern.

## Go's Component Testing Pattern

Instead of Python's `if __name__ == "__main__"`, Go uses separate executables for component testing:

- **Main executable**: `main.go` - imports and uses components
- **Component test executable**: `cmd/testboard/main.go` - tests board component independently

This follows Go's idiomatic approach where you can have multiple executables in a single project.

