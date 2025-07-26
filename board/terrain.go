package board

import (
	"fmt"
)

// TerrainType represents the different types of terrain in Catan
type TerrainType int

const (
	TerrainWood TerrainType = iota
	TerrainBrick
	TerrainOre
	TerrainWheat
	TerrainSheep
	TerrainDesert
)

// String returns the string representation of the terrain type
func (t TerrainType) String() string {
	switch t {
	case TerrainWood:
		return "Bosque"
	case TerrainBrick:
		return "Arcilla"
	case TerrainOre:
		return "Montaña"
	case TerrainWheat:
		return "Plantación"
	case TerrainSheep:
		return "Pasto"
	case TerrainDesert:
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

// RenderTile returns a 5-element array of strings representing the tile
// looks like this:
//
//		/‾‾‾‾‾‾\
//		/  B  \
//	    2
//		\      /
//		\______/
//
// where B is the terrain abbreviation and 2 is the dice number
func (tile *Tile) RenderTile(isCursor bool) [5]string {
	terrainAbbrev := tile.getTerrainAbbrev()

	diceStr := ""
	if tile.DiceNumber > 0 {
		diceStr = fmt.Sprintf("%d", tile.DiceNumber)
	}

	code := tile.getTerrainColor()
	if isCursor {
		code = fmt.Sprintf("%s\033[48;5;15m", code)
	}
	endCode := "\033[0m"
	lines := [5]string{
		fmt.Sprintf("%s/‾‾‾‾‾‾\\%s", code, endCode),
		fmt.Sprintf("%s/  %s  \\%s", code, terrainAbbrev, endCode),
		fmt.Sprintf("%s    %2s    %s", code, diceStr, endCode),
		fmt.Sprintf("%s\\        /%s", code, endCode),
		fmt.Sprintf("%s\\______/%s", code, endCode),
	}
	return lines
}

// getTerrainAbbrev returns a short abbreviation for the terrain type
func (tile *Tile) getTerrainAbbrev() string {
	switch tile.Terrain {
	case TerrainWood:
		return "BOSQ"
	case TerrainBrick:
		return "ARCI"
	case TerrainOre:
		return "MONT"
	case TerrainWheat:
		return "PLAN"
	case TerrainSheep:
		return "PAST"
	case TerrainDesert:
		return "DESI"
	default:
		return "????"
	}
}

// getTerrainColor returns the terminal 256 color code for the terrain color
func (tile *Tile) getTerrainColor() string {
	var colorNumber int
	switch tile.Terrain {
	case TerrainWood:
		// dark green
		colorNumber = 2
	case TerrainBrick:
		// red
		colorNumber = 1
	case TerrainOre:
		// dark gray
		colorNumber = 8
	case TerrainWheat:
		// yellow
		colorNumber = 11
	case TerrainSheep:
		// green
		colorNumber = 10
	case TerrainDesert:
		// brown
		colorNumber = 3
	default:
		colorNumber = 0
	}
	return fmt.Sprintf("\033[38;5;%dm", colorNumber)
}

type ResourceType string

const (
	ResourceOre   ResourceType = "Piedra"
	ResourceWood  ResourceType = "Madera"
	ResourceSheep ResourceType = "Lana"
	ResourceWheat ResourceType = "Trigo"
	ResourceBrick ResourceType = "Ladrillo"
)

func (r ResourceType) String() string {
	return string(r)
}

func TileResource(t Tile) (ResourceType, bool) {
	switch t.Terrain {
	case TerrainOre:
		return ResourceOre, true
	case TerrainWood:
		return ResourceWood, true
	case TerrainSheep:
		return ResourceSheep, true
	case TerrainWheat:
		return ResourceWheat, true
	case TerrainBrick:
		return ResourceBrick, true
	default:
		return "", false
	}
}
