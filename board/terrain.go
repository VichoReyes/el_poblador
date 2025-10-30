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
		cursorBg := lipgloss.AdaptiveColor{Light: "#E0E0E0", Dark: "#424242"}
		style = style.Background(cursorBg)
	}

	infoLine := "\\        /"
	if hasRobber {
		infoLine = "\\  ROB   /"
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
	var color lipgloss.AdaptiveColor
	switch tile.Terrain {
	case TerrainWood:
		// Dark green for dark backgrounds, forest green for light
		color = lipgloss.AdaptiveColor{Light: "#2D5016", Dark: "#4CAF50"}
	case TerrainBrick:
		// Brick red
		color = lipgloss.AdaptiveColor{Light: "#B71C1C", Dark: "#EF5350"}
	case TerrainOre:
		// Gray stone
		color = lipgloss.AdaptiveColor{Light: "#424242", Dark: "#9E9E9E"}
	case TerrainWheat:
		// Golden yellow
		color = lipgloss.AdaptiveColor{Light: "#F57F17", Dark: "#FFEB3B"}
	case TerrainSheep:
		// Bright green pasture
		color = lipgloss.AdaptiveColor{Light: "#388E3C", Dark: "#8BC34A"}
	case TerrainDesert:
		// Sandy brown
		color = lipgloss.AdaptiveColor{Light: "#795548", Dark: "#BCAAA4"}
	default:
		color = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"}
	}
	return lipgloss.NewStyle().Foreground(color)
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
