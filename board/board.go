package board

import (
	"fmt"
)

// TerrainType represents the different types of terrain in Catan
type TerrainType int

const (
	Bosque TerrainType = iota
	Arcilla
	Montaña
	Plantación
	Pasto
	Desierto
)

// String returns the string representation of the terrain type
func (t TerrainType) String() string {
	switch t {
	case Bosque:
		return "Bosque"
	case Arcilla:
		return "Arcilla"
	case Montaña:
		return "Montaña"
	case Plantación:
		return "Plantación"
	case Pasto:
		return "Pasto"
	case Desierto:
		return "Desierto"
	default:
		return "Desconocido"
	}
}

// Tile represents a hexagonal tile on the Catan board
type Tile struct {
	Terrain    TerrainType
	DiceNumber int // 2-12, 0 for desert (no dice)
}

// RenderTile returns a 4-element array of strings representing the tile
// looks like this:
//
//	/‾‾‾‾‾‾\
//	/  B  \
//	\  2   /
//	\______/
//
// where B is the terrain abbreviation and 2 is the dice number
func (tile *Tile) RenderTile() [4]string {
	terrainAbbrev := tile.getTerrainAbbrev()

	// Format dice number, empty for desert
	diceStr := ""
	if tile.DiceNumber > 0 {
		diceStr = fmt.Sprintf("%d", tile.DiceNumber)
	}

	// Create the 4 lines of the hexagon
	code := tile.getTerrainColor()
	endCode := "\033[0m"
	lines := [4]string{
		fmt.Sprintf("%s/‾‾‾‾‾\\%s", code, endCode),
		fmt.Sprintf("%s/  %s  \\%s", code, terrainAbbrev, endCode),
		fmt.Sprintf("%s\\  %2s   /%s", code, diceStr, endCode),
		fmt.Sprintf("%s\\_____/%s", code, endCode),
	}
	return lines
}

// getTerrainAbbrev returns a short abbreviation for the terrain type
func (tile *Tile) getTerrainAbbrev() string {
	switch tile.Terrain {
	case Bosque:
		return "BOS"
	case Arcilla:
		return "ARC"
	case Montaña:
		return "MTN"
	case Plantación:
		return "PLT"
	case Pasto:
		return "PAS"
	case Desierto:
		return "DES"
	default:
		return "???"
	}
}

// getTerrainColor returns the terminal 256 color code for the terrain color
func (tile *Tile) getTerrainColor() string {
	var colorNumber int
	switch tile.Terrain {
	case Bosque:
		// dark green
		colorNumber = 2
	case Arcilla:
		// red
		colorNumber = 1
	case Montaña:
		// dark gray
		colorNumber = 8
	case Plantación:
		// yellow
		colorNumber = 11
	case Pasto:
		// green
		colorNumber = 10
	case Desierto:
		// brown
		colorNumber = 3
	default:
		colorNumber = 0
	}
	return fmt.Sprintf("\033[38;5;%dm", colorNumber)
}

// PrintBoard prints the game board made of ASCII hexagons
func PrintBoard() {
	// TODO: Implement proper hexagon grid layout
	fmt.Println("TODO: Implement hexagon grid layout")
}

// CrossCoord represents the coordinates of an intersection point
// where three hexagons meet. Uses a tilted x-y coordinate system.
type CrossCoord struct {
	X, Y int
}

// TileCoord represents the coordinates of a hexagonal tile.
// Uses the same tilted x-y coordinate system but represents the tile's position.
type TileCoord struct {
	X, Y int
}

// PathCoord represents a path/edge between two intersections.
type PathCoord struct {
	From, To CrossCoord
}

// NewCrossCoord creates a new intersection coordinate
func NewCrossCoord(x, y int) CrossCoord {
	return CrossCoord{X: x, Y: y}
}

// NewTileCoord creates a new tile coordinate
func NewTileCoord(x, y int) TileCoord {
	return TileCoord{X: x, Y: y}
}

// NewPathCoord creates a new path coordinate between two intersections
func NewPathCoord(from, to CrossCoord) PathCoord {
	// Ensure canonical ordering (ascending)
	if (from.X < to.X) || (from.X == to.X && from.Y < to.Y) {
		return PathCoord{From: from, To: to}
	}
	return PathCoord{From: to, To: from}
}

// String returns the string representation of an intersection coordinate
func (ic CrossCoord) String() string {
	return fmt.Sprintf("(%d,%d)", ic.X, ic.Y)
}

// String returns the string representation of a tile coordinate
func (tc TileCoord) String() string {
	return fmt.Sprintf("[%d,%d]", tc.X, tc.Y)
}

// String returns the string representation of a path coordinate
func (pc PathCoord) String() string {
	return fmt.Sprintf("%s-%s", pc.From.String(), pc.To.String())
}
