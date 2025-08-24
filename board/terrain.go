package board

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
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
//		\  R  /
//		\______/
//
// where B is the terrain abbreviation, 2 is the dice number, and R is the robber indicator
func (tile *Tile) RenderTile(isCursor bool, hasRobber bool) [5]string {
	terrainAbbrev := tile.getTerrainAbbrev()

	diceStr := ""
	if tile.DiceNumber > 0 {
		diceStr = fmt.Sprintf("%d", tile.DiceNumber)
	}

	style := tile.getTerrainStyle()
	if isCursor {
		style = style.Background(lipgloss.Color("15"))
	}

	infoLine := "\\        /"
	if hasRobber {
		infoLine = "\\  ROB  /"
	}

	lines := [5]string{
		style.Render("/‾‾‾‾‾‾\\"),
		style.Render(fmt.Sprintf("/  %s  \\", terrainAbbrev)),
		style.Render(fmt.Sprintf("    %2s    ", diceStr)),
		style.Render(infoLine),
		style.Render("\\______/"),
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

// getTerrainStyle returns the lipgloss style for the terrain
func (tile *Tile) getTerrainStyle() lipgloss.Style {
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
	return lipgloss.NewStyle().Foreground(lipgloss.Color(fmt.Sprintf("%d", colorNumber)))
}

type ResourceType string

const (
	ResourceOre   ResourceType = "Piedra"
	ResourceWood  ResourceType = "Madera"
	ResourceSheep ResourceType = "Lana"
	ResourceWheat ResourceType = "Trigo"
	ResourceBrick ResourceType = "Ladrillo"
)

var RESOURCE_TYPES = []ResourceType{ResourceOre, ResourceWood, ResourceSheep, ResourceWheat, ResourceBrick}

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
