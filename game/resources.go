package game

import "el_poblador/board"

type ResourceType string

const (
	ResourceOre   ResourceType = "Ore"
	ResourceWood  ResourceType = "Wood"
	ResourceSheep ResourceType = "Sheep"
	ResourceWheat ResourceType = "Wheat"
	ResourceBrick ResourceType = "Brick"
)

func (r ResourceType) String() string {
	return string(r)
}

func tileResource(t board.Tile) (ResourceType, bool) {
	switch t.Terrain {
	case board.Montaña:
		return ResourceOre, true
	case board.Bosque:
		return ResourceWood, true
	case board.Pasto:
		return ResourceSheep, true
	case board.Plantación:
		return ResourceWheat, true
	case board.Arcilla:
		return ResourceBrick, true
	default:
		return "", false
	}
}
