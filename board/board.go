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

// GetBoardSize returns the dimensions of the board
func GetBoardSize() (int, int) {
	// TODO: Return actual board dimensions
	return 5, 5 // placeholder
}
